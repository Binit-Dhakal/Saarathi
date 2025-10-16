package domain

const (
	TripCreatedEvent = "trips.created"
)

type TripCreated struct {
	SagaID   string
	TripID   string
	Pickup   [2]float64 // Lng, Lat
	DropOff  [2]float64
	Distance float64
	Price    int
	CarType  CarPackage
}
