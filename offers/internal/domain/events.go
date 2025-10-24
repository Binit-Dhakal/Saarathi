package domain

const (
	RideMatchingInitializedEvent = "offers.rms.initialized"
	TripOfferEvent               = "offers.drivers.request"
	TripOfferAcceptedEvent       = "offers.request.accepted"
	NoCandidateMatchedEvent      = "offers.rms.notFound"
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
	PickUp           [2]float64
	DropOff          [2]float64
	DriverID         string
	PresenceServerID string
}

type TripOfferAccepted struct {
	SagaID   string
	TripID   string
	DriverID string
}

type NoCandidateMatched struct {
	SagaID            string
	TripID            string
	PickUp            [2]float64
	DropOff           [2]float64
	CarType           string
	MaxSearchRadiusKm int32
	Attempt           int32
	FirstAttemptUnix  int64
}
