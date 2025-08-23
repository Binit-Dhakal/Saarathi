package domain

// driver location
type Location struct {
	Lat float64
	Lon float64
}

type DriverLocation struct {
	DriverID    string
	Longitude   float64
	Latitude    float64
	VehicleType string
}
