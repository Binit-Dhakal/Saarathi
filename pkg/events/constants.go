package events

var (
	EventTripCreated   = TripEventCreated{}.EventName()
	EventOfferRequest  = TripOfferRequest{}.EventName()
	EventOfferResponse = TripOfferResponse{}.EventName()
)
