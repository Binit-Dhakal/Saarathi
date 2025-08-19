package domain

type RideStatus string

var (
	RideStatusPending   RideStatus = "pending"
	RideStatusApproved  RideStatus = "approved"
	RideStatusCancelled RideStatus = "cancelled"
)

type RideModel struct {
	RideID   string
	RiderID  string
	DriverID string
	FareID   string
	Status   RideStatus
}
