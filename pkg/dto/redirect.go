package dto

// RedirectResponse is a huma output model for handlers that respond with an
// HTTP redirect. Status is read by huma to set the response code dynamically.
type RedirectResponse struct {
	Status   int
	Location string `header:"Location"`
}
