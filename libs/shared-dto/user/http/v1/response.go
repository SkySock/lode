package v1

type SignUpResponse struct {
	UserId string `json:"userId"`
}

type SignInResponse struct {
	AccessToken string `json:"accessToken"`
}
