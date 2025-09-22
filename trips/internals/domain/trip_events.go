package domain

import "time"

type TripCreated struct {
	TripID     string
	RiderID    string
	Pickup     Coordinate
	DropOff    Coordinate
	OccurredAt time.Time
}
