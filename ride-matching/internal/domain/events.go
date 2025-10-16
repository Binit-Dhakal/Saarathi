package domain

const MatchingCandidatesEvent = "rms.matching.candidates"

type MatchingCandidates struct {
	SagaID    string
	TripID    string
	DriverIds []string
}
