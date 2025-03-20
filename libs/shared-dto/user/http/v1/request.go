package v1

import (
	"encoding/json"
	"io"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type SignUpRequest struct {
	Username string `json:"username" validate:"required,alphanum,gte=1,lte=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
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
	validator.RegisterValidation("password", validatePassword)

	return validator.Struct(s)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 {
		return false
	}

	// Upper, Lower, Number, Special flags
	ulns := 0
	specialChars := "!@#$%^&*._"

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			ulns |= 1 << 3
		case unicode.IsLower(char):
			ulns |= 1 << 2
		case unicode.IsDigit(char):
			ulns |= 1 << 1
		case strings.ContainsRune(specialChars, char):
			ulns |= 1
		case unicode.IsSpace(char):
			return false
		}
	}

	if ulns&0b1111 != 0b1111 {
		return false
	}

	return true
}
