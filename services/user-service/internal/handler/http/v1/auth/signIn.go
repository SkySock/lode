package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	v1 "github.com/SkySock/libs/shared-dto/user/http/v1"
	"github.com/SkySock/lode/libs/utils/http/response"
	"github.com/SkySock/lode/services/user-service/internal/usecase"
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

func (h *SignIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.l.Warn("invalid request method", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := v1.SignInRequest{}

	if err := data.FromJSON(r.Body); err != nil {
		http.Error(w, "Error input data", http.StatusBadRequest)
		return
	}

	data.Normalize()

	if err := data.Validate(); err != nil {
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
