package domain

import "github.com/Binit-Dhakal/Saarathi/pkg/ddd"

type RideStatus string

var (
	RideStatusPending   RideStatus = "pending"
	RideStatusApproved  RideStatus = "approved"
	RideStatusCancelled RideStatus = "cancelled"
)

type RideModel struct {
	ddd.Aggregate
	RideID   string
	RiderID  string
	DriverID string
	FareID   string
	Status   RideStatus
}
