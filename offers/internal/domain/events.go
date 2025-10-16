package domain

const (
	RideMatchingInitializedEvent = "offers.rms.initialized"
	TripOfferEvent               = "offers.drivers.request"
)

type RideMatchingInitialized struct {
	SagaID  string
	TripID  string
	PickUp  [2]float64
	DropOff [2]float64
	CarType string
}

type TripOffer struct {
	SagaID           string
	TripID           string
	Price            int32
	Distance         float64
	PresenceServerID string
}
