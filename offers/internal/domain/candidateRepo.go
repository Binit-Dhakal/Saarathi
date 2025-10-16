package domain

import "context"

type TripCandidatesRepository interface {
	SaveCandidates(ctx context.Context, tripID string, driverIDs []string) error
	IncrementCandidateCounter(ctx context.Context, tripID string) (int, error)
	GetNextCandidates(ctx context.Context, tripID string, index int) (string, error)
}
