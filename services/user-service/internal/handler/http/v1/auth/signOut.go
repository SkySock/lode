package auth

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type SignOut struct {
	l  *slog.Logger
	uc logoutUsecase
}

type logoutUsecase interface {
	Logout(ctx context.Context, refreshToken string) error
}

func NewSignOut(l *slog.Logger, uc logoutUsecase) *SignOut {
	return &SignOut{l, uc}
}

var _ http.Handler = (*SignOut)(nil)

// SignOut godoc
// @Summary      Выход из системы
// @Description  Делает недействительным refresh токен и удаляет файл cookie refreshToken
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Success      200 {string} string "Successfully logged out"
// @Failure      400 {string} string "Missing or empty refreshToken cookie"
// @Failure      500 {string} string "Failed to logout"
// @Router       /auth/sign-out [post]
func (h *SignOut) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.l.Warn("invalid request method", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	refresh, err := r.Cookie("refreshToken")
	if err != nil || refresh.Value == "" {
		h.l.Warn("missing or empty refreshToken cookie", "error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err = h.uc.Logout(r.Context(), refresh.Value); err != nil {
		h.l.Error("failed to logout", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   true,
	})
	w.WriteHeader(http.StatusOK)
}
