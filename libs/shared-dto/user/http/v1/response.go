package v1

type SignUpResponse struct {
	UserId string `json:"userId" example:"01976451-00b3-7e32-9340-4f999c6c5edd"`
}

type SignInResponse struct {
	AccessToken string `json:"accessToken"`
}
