package v1

import (
	"encoding/json"
	"io"
)

type SignUpResponse struct {
	UserId string `json:"userId"`
}

type SignInResponse struct {
	AccessToken string `json:"accessToken"`
}

func (s *SignUpResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(s)
}
