package application

import (
	"context"
	"fmt"
	"strings"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/dto"
)

type MatchingService interface {
	ProcessMatchingRequest(ctx context.Context, requestDTO dto.TripCreated) error
}

type matchingService struct {
	publisher   ddd.EventPublisher[ddd.Event]
	availCheck  domain.DriverAvailabilityChecker
	metaFetcher domain.DriverMetadataFetcher
	matchRepo   domain.RedisRideMatchingRepository
}

func NewMatchingService(publisher ddd.EventPublisher[ddd.Event], matchRepo domain.RedisRideMatchingRepository, metaFetcher domain.DriverMetadataFetcher, availCheck domain.DriverAvailabilityChecker) MatchingService {
	return &matchingService{
		publisher:   publisher,
		matchRepo:   matchRepo,
		metaFetcher: metaFetcher,
		availCheck:  availCheck,
	}
}

// currently our algorithm just find drivers based on geographical location
func (m *matchingService) ProcessMatchingRequest(ctx context.Context, requestDTO dto.TripCreated) error {
	candidates := m.matchRepo.FindNearestDriver(ctx, requestDTO.PickUp.Lng, requestDTO.PickUp.Lat)

	// availability check
	onlineCandidates := m.availCheck.GetOnlineDrivers(ctx, candidates)

	if len(onlineCandidates) == 0 {
		return fmt.Errorf("No driver online")
	}

	metadatas, err := m.metaFetcher.GetBulkMetada(ctx, onlineCandidates)
	if err != nil {
		return fmt.Errorf("failed to fetch driver metadata: %w", err)
	}

	var shortlistDrivers []string
	selectedVehicle := strings.ToUpper(requestDTO.CarType)

	for _, metadata := range metadatas {
		vt := strings.ToUpper(metadata.VehicleType)
		if vt == selectedVehicle {
			shortlistDrivers = append(shortlistDrivers, metadata.DriverID)
		}
	}

	replyPayload := &domain.MatchingCandidates{
		SagaID:           requestDTO.SagaID,
		TripID:           requestDTO.TripID,
		DriverIds:        shortlistDrivers,
		Attempt:          requestDTO.Attempt,
		SearchRadius:     requestDTO.SearchRadius,
		FirstAttemptUnix: requestDTO.FirstAttemptUnix,
	}
	matchEvt := ddd.NewEvent(domain.MatchingCandidatesEvent, replyPayload)

	return m.publisher.Publish(ctx, matchEvt)
}
