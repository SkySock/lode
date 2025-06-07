package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/SkySock/lode/services/user-service/internal/usecase"
	"github.com/SkySock/lode/services/user-service/internal/validation"

	v1 "github.com/SkySock/lode/libs/shared-dto/user/http/v1"
	"github.com/SkySock/lode/libs/utils/http/response"
	"github.com/google/uuid"
)

type registerUsecase interface {
	RegisterUser(ctx context.Context, userInfo usecase.RegistrationInfo) (uuid.UUID, error)
}

type SignUp struct {
	l  *slog.Logger
	uc registerUsecase
}

func NewSignUp(l *slog.Logger, uc registerUsecase) *SignUp {
	return &SignUp{l, uc}
}

var _ http.Handler = (*SignUp)(nil)

func (h *SignUp) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.l.Warn("invalid request method", "method", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	data := v1.SignUpRequest{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Error input data", http.StatusBadRequest)
		return
	}
	data.Normalize()

	err = validation.ValidateSignUpRequest(&data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error validating data: %s", err), http.StatusBadRequest)
		return
	}

	newUser := usecase.RegistrationInfo{
		Username: data.Username,
		Email:    data.Email,
		Password: data.Password,
	}

	userId, err := h.uc.RegisterUser(r.Context(), newUser)
	if err != nil {
		if errors.Is(err, usecase.ErrEmailOrUsernameAlreadyExists) {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	resp := &v1.SignUpResponse{
		UserId: userId.String(),
	}

	if err := response.WriteJSON(w, http.StatusCreated, resp); err != nil {
		h.l.Error("JSON encoding failed", "error", err)
		http.Error(w, "Unable encode json", http.StatusInternalServerError)
		return
	}
}
