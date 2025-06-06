package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/SkySock/lode/services/user-service/internal/entity"
	"github.com/valkey-io/valkey-go"
)

// TODO: implements cleenup sessions

const (
	tokenPrefix   = "refresh_token"
	sessionPrefix = "user_session"
)

type sessionRepository struct {
	client valkey.Client
}

func NewSessionRepository(client valkey.Client) SessionRepository {
	return &sessionRepository{
		client: client,
	}
}

func (r *sessionRepository) Save(ctx context.Context, token string, data *entity.Session) error {
	tokenKey := fmt.Sprintf("%s:%s", tokenPrefix, token)
	sessionKey := fmt.Sprintf("%s:%s", sessionPrefix, data.UserID)

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	ttl := time.Until(data.ExpiresAt)

	c, cancel := r.client.Dedicate()
	defer cancel()

	c.Do(ctx, c.B().Watch().Key(tokenKey, sessionKey).Build())
	cmds := valkey.Commands{
		c.B().Multi().Build(),
		c.B().Set().Key(tokenKey).Value(string(bytes)).Ex(ttl).Build(),
		c.B().Sadd().Key(fmt.Sprintf("%s:%s", sessionPrefix, data.UserID)).Member(token).Build(),
		c.B().Exec().Build(),
	}

	results := c.DoMulti(ctx, cmds...)

	for i, res := range results {
		if err := res.Error(); err != nil {
			if i == len(results)-1 {
				if valkey.IsValkeyNil(err) {
					return errors.New("transaction aborted due to concurrent modification")
				}
			}
			return err
		}
	}

	return nil
}

func (r *sessionRepository) Get(ctx context.Context, token string) (*entity.Session, error) {
	key := fmt.Sprintf("%s:%s", tokenPrefix, token)
	cmd := r.client.B().Get().Key(key).Build()

	result := r.client.Do(ctx, cmd)
	if err := result.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	reader, err := result.AsReader()
	if err != nil {
		return nil, err
	}

	var sessionData entity.Session

	err = json.NewDecoder(reader).Decode(&sessionData)
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (r *sessionRepository) Revoke(ctx context.Context, token string) error {
	key := fmt.Sprintf("%s:%s", tokenPrefix, token)

	script := `
	local key = KEYS[1]
	local json = redis.call('GET', key)
	if not json then return nil end

	local data = cjson.decode(json)
	data['revoked'] = true
	redis.call('SET', key, cjson.encode(data), 'KEEPTTL')
	return 1
	`

	vscript := valkey.NewLuaScript(script)
	resp := vscript.Exec(ctx, r.client, []string{key}, []string{})
	if err := resp.Error(); err != nil {
		if valkey.IsValkeyNil(err) {
			return ErrNotFound
		}
		return err
	}
	return nil
}

func (r *sessionRepository) Delete(ctx context.Context, token string) error {
	session, err := r.Get(ctx, token)
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return nil
		}
		return err
	}

	tokenKey := fmt.Sprintf("%s:%s", tokenPrefix, token)
	sessionKey := fmt.Sprintf("%s:%s", sessionPrefix, session.UserID)

	c, cancel := r.client.Dedicate()
	defer cancel()

	watchResp := c.Do(ctx, c.B().Watch().Key(tokenKey, sessionKey).Build())
	if err := watchResp.Error(); err != nil {
		return err
	}

	cmds := valkey.Commands{
		c.B().Multi().Build(),
		c.B().Del().Key(tokenKey).Build(),
		c.B().Srem().Key(sessionKey).Member(token).Build(),
		c.B().Exec().Build(),
	}

	results := c.DoMulti(ctx, cmds...)

	for i, res := range results {
		if err := res.Error(); err != nil {
			if i == len(results)-1 && valkey.IsValkeyNil(err) {
				return errors.New("transaction aborted due to concurrent modification")
			}
			return err
		}
	}

	return nil
}

func (r *sessionRepository) GetUserSessions(ctx context.Context, userID string) ([]string, error) {
	key := fmt.Sprintf("%s:%s", sessionPrefix, userID)
	cmd := r.client.B().Smembers().Key(key).Build()

	result := r.client.Do(ctx, cmd)
	if err := result.Error(); err != nil {
		return nil, err
	}

	return result.AsStrSlice()
}
