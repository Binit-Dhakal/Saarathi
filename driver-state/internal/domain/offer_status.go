package domain

type OfferStatus string

const (
	OfferPending   OfferStatus = "pending"
	OfferDelivered OfferStatus = "delivered"
	OfferAccepted  OfferStatus = "accepted"
	OfferRejected  OfferStatus = "rejected"
	OfferTimedOut  OfferStatus = "timed_out"
	OfferCancelled OfferStatus = "cancelled"
)
