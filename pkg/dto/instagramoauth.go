package dto

type InstagramStartRequest struct{}

type InstagramCallbackRequest struct {
	Code string `query:"code" required:"true"`
}
