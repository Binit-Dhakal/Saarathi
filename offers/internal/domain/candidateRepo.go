package domain

import "context"

type TripCandidatesRepository interface {
	SaveCandidates(ctx context.Context, tripID string, driverIDs []string) error
	IncrementCandidateCounter(ctx context.Context, tripID string) (int, error)
	GetNextCandidates(ctx context.Context, tripID string, index int) (string, error)
	SaveFirstAttemptUnix(ctx context.Context, tripID string) error
	GetFirstAttemptUnix(ctx context.Context, tripID string) (int64, error)
	AddRejectedDriver(ctx context.Context, tripID string, driverID string) error
}
