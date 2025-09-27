package domain

import "time"

// driver location
type Location struct {
	Lat       float64
	Lng       float64
	UpdatedAt time.Time
}

type DriverLocation struct {
	DriverID    string
	Longitude   float64
	Latitude    float64
	VehicleType string
}
