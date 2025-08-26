package events

var (
	EventTripCreated   = TripEventCreated{}.EventName()
	EventOfferCreated  = TripOffer{}.EventName()
	EventOfferResponse = TripOfferResponse{}.EventName()
)
