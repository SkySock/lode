package v1

import (
	"strings"
)

type SignUpRequest struct {
	Username string `json:"username" example:"ozon671games" validate:"required,alphanum,gte=1,lte=50"`
	Email    string `json:"email" example:"example@example.com" validate:"required,email"`
	Password string `json:"password" example:"Da1dfshgn$" validate:"required,password"`
}

type SignInRequest struct {
	Login    string `json:"login" example:"ozon671games" validate:"required"`
	Password string `json:"password" example:"Da1dfshgn$" validate:"required"`
}

func (s *SignUpRequest) Normalize() {
	s.Username = strings.ToLower(s.Username)
	s.Email = strings.ToLower(s.Email)
}

func (s *SignInRequest) Normalize() {
	s.Login = strings.ToLower(s.Login)
}
