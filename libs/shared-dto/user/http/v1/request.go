package v1

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
)

type SignUpRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (s *SignUpRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(s)
}

func (s *SignUpRequest) Validate() error {
	validator := validator.New()
	return validator.Struct(s)
}
