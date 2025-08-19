package application

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
)

type RouteService interface {
	GetRouteDetailFromOSRM(domain.Coordinate, domain.Coordinate) (*domain.Route, error)
}

type routeService struct {
	client *http.Client
}

func NewRouteService() *routeService {
	return &routeService{
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (r *routeService) GetRouteDetailFromOSRM(pickOff domain.Coordinate, dropOff domain.Coordinate) (*domain.Route, error) {
	baseURL, err := url.Parse(
		fmt.Sprintf("http://osrm-backend:5000/route/v1/driving/%v,%v;%v,%v",
			pickOff.Lon, pickOff.Lat, dropOff.Lon, dropOff.Lat,
		),
	)
	if err != nil {
		return nil, err
	}
	queryParams := baseURL.Query()
	queryParams.Add("overview", "full")
	queryParams.Add("alternatives", "false")
	queryParams.Add("annotations", "false")
	queryParams.Add("steps", "false")
	queryParams.Add("geometries", "geojson")

	baseURL.RawQuery = queryParams.Encode()

	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OSRM error: %s", body)
	}

	var osrmResponse domain.OSRMResponse
	if err := json.NewDecoder(resp.Body).Decode(&osrmResponse); err != nil {
		return nil, err
	}

	var routeRes domain.Route

	routeRes = osrmResponse.Route[0]
	routeRes.Source = pickOff
	routeRes.Destination = dropOff

	return &routeRes, nil
}
