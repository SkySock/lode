package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	v1 "github.com/SkySock/lode/libs/shared-dto/user/http/v1"
	"github.com/SkySock/lode/libs/utils/http/response"
	"github.com/SkySock/lode/services/user-service/internal/usecase"
	"github.com/SkySock/lode/services/user-service/internal/validation"
)

type SignIn struct {
	l  *slog.Logger
	uc loginUsecase
}

type loginUsecase interface {
	Login(ctx context.Context, login, password string) (*usecase.AuthTokens, error)
}

func NewSignIn(l *slog.Logger, uc loginUsecase) *SignIn {
	return &SignIn{l, uc}
}

var _ http.Handler = (*SignIn)(nil)

// SignIn godoc
// @Summary      Вход в аккаунт
// @Description  Вход в аккаунт пользователя
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body v1.SignInRequest true "Данные регистрации"
// @Success      200  {object}  v1.SignInResponse
// @Router       /auth/sign-in [post]
func (h *SignIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.l.Warn("invalid request method", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := v1.SignInRequest{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Error input data", http.StatusBadRequest)
		return
	}

	data.Normalize()

	if err := validation.ValidateSignInRequest(&data); err != nil {
		http.Error(w, fmt.Sprintf("Error validating data: %s", err), http.StatusBadRequest)
		return
	}

	tokens, err := h.uc.Login(r.Context(), data.Login, data.Password)
	if err != nil {
		if errors.Is(err, usecase.ErrIncorrectPassword) ||
			errors.Is(err, usecase.ErrEmailNotFound) ||
			errors.Is(err, usecase.ErrUsernameNotFound) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		h.l.Error("Login failed", "error", err)
		http.Error(w, "Authentication failed", http.StatusInternalServerError)
		return
	}

	refreshToken := http.Cookie{
		Name:     "refreshToken",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, &refreshToken)

	responseBody := &v1.SignInResponse{
		AccessToken: tokens.AccessToken,
	}

	if err := response.WriteJSON(w, http.StatusOK, responseBody); err != nil {
		h.l.Error("JSON encoding failed", "error", err)
		http.Error(w, "Unable encode json", http.StatusInternalServerError)
		return
	}
}
