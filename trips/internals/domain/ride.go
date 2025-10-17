package domain

import "github.com/Binit-Dhakal/Saarathi/pkg/ddd"

type TripModel struct {
	ddd.Aggregate
	TripID   string
	RiderID  string
	DriverID string
	FareID   string
	Status   TripStatus
}

type TripProjectionDetail struct {
	TripID    string
	FarePrice int
	Distance  float64
	PickUp    Coordinate
	DropOff   Coordinate
}
