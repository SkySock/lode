package validation

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 8 || strings.ContainsAny(password, " \t\n\r") {
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
