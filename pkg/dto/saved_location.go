package dto

type SavedLocationResponse struct {
	ID     string  `json:"id"`
	UserID string  `json:"userId"`
	Name   string  `json:"name"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
}

type SavedLocationsResponse struct {
	SavedLocations []SavedLocationResponse `json:"savedLocations"`
}

type CreateSavedLocationRequest struct {
	Body struct {
		Name string  `json:"name" required:"true"`
		Lat  float64 `json:"lat" required:"true"`
		Lng  float64 `json:"lng" required:"true"`
	}
}

type ListSavedLocationsRequest struct{}

type GetSavedLocationRequest struct {
	ID string `path:"id" required:"true"`
}

type UpdateSavedLocationRequest struct {
	ID   string `path:"id" required:"true"`
	Body struct {
		Name string  `json:"name" required:"true"`
		Lat  float64 `json:"lat" required:"true"`
		Lng  float64 `json:"lng" required:"true"`
	}
}

type DeleteSavedLocationRequest struct {
	ID string `path:"id" required:"true"`
}

type DeleteSavedLocationResponse struct{}
