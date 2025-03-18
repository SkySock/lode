package v1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SkySock/lode/services/user-service/internal/usecase"

	v1 "github.com/SkySock/lode/libs/share-dto/user/http/v1"
	"github.com/google/uuid"
)

type userUsecase interface {
	RegisterUser(ctx context.Context, userInfo usecase.RegistrationInfo) (uuid.UUID, error)
}

type SignUp struct {
	l           *slog.Logger
	userUsecase userUsecase
}

func NewSignUp(l *slog.Logger, userUsecase userUsecase) *SignUp {
	return &SignUp{l, userUsecase}
}

var _ http.Handler = (*SignUp)(nil)

func (s *SignUp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data := v1.SignUpRequest{}

	err := data.FromJSON(r.Body)
	if err != nil {
		http.Error(w, "Error input data", http.StatusBadRequest)
		return
	}

	err = data.Validate()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error validating data: %s", err), http.StatusBadRequest)
		return
	}

	newUser := usecase.RegistrationInfo{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}

	userId, err := s.userUsecase.RegisterUser(r.Context(), newUser)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailOrUsernameAlreadyExists) {
			http.Error(w, "Username or email already exists", http.StatusBadRequest)
			return
		}
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	resp := &v1.SignUpResponse{
		UserId: userId.String(),
	}
	err = resp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable encode json", http.StatusInternalServerError)
		return
	}
}
