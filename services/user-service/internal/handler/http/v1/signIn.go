package v1

import (
	"log/slog"
	"net/http"
)

type SignIn struct {
	l             *slog.Logger
	signInUsecase signInUsecase
}

type signInUsecase interface{}

func NewSignIn(l *slog.Logger, usecase signInUsecase) *SignIn {
	return &SignIn{l, usecase}
}

var _ http.Handler = (*SignIn)(nil)

func (s *SignIn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// TODO
}
