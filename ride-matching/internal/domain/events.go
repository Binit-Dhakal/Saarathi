package domain

const (
	MatchingCandidatesEvent = "rms.matching.candidates"
	NoDriverAvailableEvent  = "rms.matching.notAvailable"
)

type MatchingCandidates struct {
	SagaID    string
	TripID    string
	DriverIds []string

	SearchRadius     int32
	Attempt          int32
	FirstAttemptUnix int64
}

type NoDriverAvailable struct {
	SagaID string
	TripID string

	Attempt          int32
	FirstAttemptUnix int64
}
