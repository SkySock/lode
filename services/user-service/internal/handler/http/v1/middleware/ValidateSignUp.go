package middleware

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/SkySock/libs/shared-dto/user/http/v1"
)

type KeySignUp struct{}

func ValidateSignUp(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		ctx := context.WithValue(r.Context(), KeySignUp{}, data)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
