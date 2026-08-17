package dto

type RefreshTokenRequest struct {
	Body struct {
		RefreshToken string `json:"refreshToken" required:"true"`
	}
}

type TokenPairResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
