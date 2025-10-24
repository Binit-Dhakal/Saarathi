package domain

type TripReadModelDTO struct {
	SagaID   string
	TripID   string
	PickUp   [2]float64 // Lng,Lat
	DropOff  [2]float64 // Lng, Lat
	CarType  string
	Price    int32
	Distance float64
}

type MatchedDriversDTO struct {
	TripID             string
	SagaID             string
	CandidateDriversID []string
	Attempt            int32
	FirstAttemptUnix   int64
}

type OfferAcceptedReplyDTO struct {
	SagaID   string
	TripID   string
	DriverID string
}
