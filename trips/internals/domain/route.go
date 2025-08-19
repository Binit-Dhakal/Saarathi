package domain

import "encoding/json"

type Coordinate struct {
	Lat float64
	Lon float64
}

type Geometry struct {
	Coordinates []Coordinate `json:"coordinates"`
	Type        string       `json:"type"`
}

type Route struct {
	Geometry    Geometry   `json:"geometry"`
	Duration    float64    `json:"duration"`
	Distance    float64    `json:"distance"`
	Source      Coordinate `json:"source"`
	Destination Coordinate `json:"destination"`
}

type OSRMResponse struct {
	Route []Route `json:"routes"`
}

func (c *Coordinate) UnmarshalJSON(data []byte) error {
	var coords [2]float64
	if err := json.Unmarshal(data, &coords); err != nil {
		return err
	}

	c.Lon = coords[0]
	c.Lat = coords[1]
	return nil
}
