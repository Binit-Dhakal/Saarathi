package domain

const (
	RideMatchingInitializedEvent = "offers.rms.initialized"
)

type RideMatchingInitialized struct {
	SagaID  string
	TripID  string
	PickUp  [2]float64
	DropOff [2]float64
	CarType string
}
