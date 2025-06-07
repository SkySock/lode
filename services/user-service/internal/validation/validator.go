package validation

import (
	v1 "github.com/SkySock/lode/libs/shared-dto/user/http/v1"
	"github.com/go-playground/validator/v10"
)

var v = func() *validator.Validate {
	val := validator.New()
	val.RegisterValidation("password", validatePassword)
	return val
}()

func ValidateSignUpRequest(body *v1.SignUpRequest) error {
	return v.Struct(body)
}

func ValidateSignInRequest(body *v1.SignInRequest) error {
	return v.Struct(body)
}
