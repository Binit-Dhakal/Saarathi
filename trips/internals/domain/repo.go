package domain

// ephemeral save of fare data in redis
type FareRepository interface {
	CreateEphemeralFareEntry(fare *FareQuote) (string, error)
	GetEphemeralFareEntry(id string) (*FareQuote, error)
	DeleteEphemeralFareEntry(id string) error
}

type TripRepository interface {
	SaveRouteDetail(route *Route, riderID string) (string, error)
	SaveFareDetail(fareModel FareRecord) (string, error)
	SaveRideDetail(rideModel TripModel) (string, error)
	AssignDriverToTrip(tripID string, driverID string) (string, error) // riderID return
}
