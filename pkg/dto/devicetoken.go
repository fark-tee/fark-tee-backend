package dto

type RegisterDeviceTokenRequest struct {
	Body struct {
		Token    string `json:"token" required:"true"`
		Platform string `json:"platform" required:"true" enum:"ANDROID,IOS"`
	}
}

type RegisterDeviceTokenResponse struct{}

type DeleteDeviceTokenRequest struct {
	Body struct {
		Token string `json:"token" required:"true"`
	}
}

type DeleteDeviceTokenResponse struct{}
