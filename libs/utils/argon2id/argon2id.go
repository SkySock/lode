package argon2id

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/argon2"
)

// Params represents the parameters for argon2id hashing
type Params struct {
	Memory      uint32 // memory in kibibytes (KiB)
	Iterations  uint32 // number of passes
	SaltLength  uint32 // salt length in bytes
	KeyLength   uint32 // key length in bytes
	Parallelism uint8  // degree of parallelism
}

type argon2Hash struct {
	algorithm   string
	version     int
	memory      uint32
	iterations  uint32
	parallelism uint8
	salt        []byte
	hash        []byte
}

func HashPassword(password []byte, params *Params) (string, error) {
	if params == nil {
		return "", errors.New("params cannot be nil")
	}
	if params.Iterations < 1 {
		return "", errors.New("iterations must be at least 1")
	}
	if params.Parallelism < 1 {
		return "", errors.New("parallelism must be at least 1")
	}
	if params.KeyLength < 4 {
		return "", errors.New("keyLength must be at least 4")
	}
	if params.SaltLength < 8 {
		return "", errors.New("saltLength must be at least 8")
	}

	minMemory := 8 * uint32(params.Parallelism)
	if params.Memory < minMemory {
		return "", fmt.Errorf(
			"memory must be at least %d KiB (8 * parallelism = %d)",
			minMemory,
			params.Parallelism,
		)
	}

	salt, err := generateRandomBytes(params.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey(
		password,
		salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash := fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version,
		params.Memory,
		params.Iterations,
		params.Parallelism,
		b64Salt,
		b64Hash,
	)

	return encodedHash, nil
}

func VerifyPassword(password []byte, encodedHash string) (bool, error) {
	parsedHash, err := parsePasswordHash(encodedHash)
	if err != nil {
		return false, fmt.Errorf("failed to parse hash: %w", err)
	}

	newHash := argon2.IDKey(
		password,
		parsedHash.salt,
		parsedHash.iterations,
		parsedHash.memory,
		parsedHash.parallelism,
		uint32(len(parsedHash.hash)),
	)

	if subtle.ConstantTimeCompare(parsedHash.hash, newHash) == 1 {
		return true, nil
	}

	return false, nil
}

func generateRandomBytes(length uint32) ([]byte, error) {
	if length == 0 {
		return nil, errors.New("salt length cannot be zero")
	}

	b := make([]byte, length)

	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func parsePasswordHash(encodedHash string) (*argon2Hash, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, errors.New("invalid hash format")
	}

	result := &argon2Hash{}

	// Parse algorithm
	if parts[1] != "argon2id" {
		return nil, fmt.Errorf("unsupported algorithm: %s", parts[1])
	}
	result.algorithm = parts[1]

	// Parse version
	versionPart := strings.TrimPrefix(parts[2], "v=")
	version, err := strconv.Atoi(versionPart)
	if err != nil {
		return nil, fmt.Errorf("invalid version: %w", err)
	}
	if version != argon2.Version {
		return nil, fmt.Errorf("unsupported argon2 version: %d", version)
	}
	result.version = version

	// Parse parameters
	params := strings.Split(parts[3], ",")
	if len(params) != 3 {
		return nil, errors.New("invalid parameters section")
	}

	for _, param := range params {
		switch {
		case strings.HasPrefix(param, "m="):
			m, err := strconv.ParseUint(strings.TrimPrefix(param, "m="), 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid memory parameter: %w", err)
			}
			result.memory = uint32(m)

		case strings.HasPrefix(param, "t="):
			t, err := strconv.ParseUint(strings.TrimPrefix(param, "t="), 10, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid iterations parameter: %w", err)
			}
			result.iterations = uint32(t)

		case strings.HasPrefix(param, "p="):
			p, err := strconv.ParseUint(strings.TrimPrefix(param, "p="), 10, 8)
			if err != nil {
				return nil, fmt.Errorf("invalid parallelism parameter: %w", err)
			}
			result.parallelism = uint8(p)

		default:
			return nil, fmt.Errorf("unknown parameter: %s", param)
		}
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return nil, fmt.Errorf("salt decoding failed: %w", err)
	}
	result.salt = salt

	// Decode hash
	hashBytes, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, fmt.Errorf("hash decoding failed: %w", err)
	}
	result.hash = hashBytes

	return result, nil
}
