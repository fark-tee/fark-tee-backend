package osrm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/fark-tee/fark-tee-backend/internal/config"
)

// Route is the OSRM-computed travel time and distance between two points.
type Route struct {
	DurationSeconds int
	DistanceMeters  float64

	// Geometry is the route's road-following path, encoded as a Google
	// polyline (precision 5). Only populated when Route is called with
	// withGeometry=true.
	Geometry string
}

type Client struct {
	cfg        *config.Config
	httpClient *http.Client
}

// @WireSet("Infrastructure")
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}
}

// Route calls the OSRM /route/v1 endpoint to estimate travel time and
// distance from (fromLat, fromLng) to (toLat, toLng) using the configured
// profile (e.g. "driving", "walking", "cycling"). When withGeometry is true,
// the returned Route also carries the road-following path as an encoded
// polyline - callers that only need duration/distance (e.g. a per-tick ETA
// refresh) should pass false to avoid the extra payload.
func (c *Client) Route(ctx context.Context, fromLat, fromLng, toLat, toLng float64, withGeometry bool) (Route, error) {
	overview := "false"
	if withGeometry {
		overview = "full"
	}

	url := fmt.Sprintf(
		"%s/route/v1/%s/%f,%f;%f,%f?overview=%s&geometries=polyline",
		strings.TrimRight(c.cfg.OSRM.BaseURL, "/"),
		c.cfg.OSRM.Profile,
		fromLng, fromLat, toLng, toLat,
		overview,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Route{}, err
	}

	var body struct {
		Code   string `json:"code"`
		Routes []struct {
			Duration float64 `json:"duration"`
			Distance float64 `json:"distance"`
			Geometry string  `json:"geometry"`
		} `json:"routes"`
	}

	if err := doJSON(c.httpClient, req, &body); err != nil {
		return Route{}, fmt.Errorf("osrm: route request failed: %w", err)
	}

	if body.Code != "Ok" || len(body.Routes) == 0 {
		return Route{}, fmt.Errorf("osrm: no route found (code=%s)", body.Code)
	}

	return Route{
		DurationSeconds: int(math.Round(body.Routes[0].Duration)),
		DistanceMeters:  body.Routes[0].Distance,
		Geometry:        body.Routes[0].Geometry,
	}, nil
}

func doJSON(client *http.Client, req *http.Request, out any) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(respBody))
	}

	return json.Unmarshal(respBody, out)
}
