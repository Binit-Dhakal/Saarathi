package application

import (
	"context"
	"fmt"
	"slices"
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

func (m *matchingService) exclude(candidates []string, rejected []string) []string {
	rejectedMap := make(map[string]bool)
	for _, item := range rejected {
		rejectedMap[item] = true
	}

	filtered := slices.DeleteFunc(candidates, func(candidate string) bool {
		return rejectedMap[candidate] // Returns true if the element should be deleted
	})

	return filtered
}

// currently our algorithm just find drivers based on geographical location
func (m *matchingService) ProcessMatchingRequest(ctx context.Context, requestDTO dto.TripCreated) error {
	const MaxRadiusKm = 5
	radius := min(requestDTO.Attempt+1, MaxRadiusKm)
	rejected, _ := m.matchRepo.GetRejectedDriver(ctx, requestDTO.TripID)
	var shortlistDrivers []string

	for radius <= MaxRadiusKm {
		candidates := m.matchRepo.FindNearestDriver(ctx, requestDTO.PickUp.Lng, requestDTO.PickUp.Lat, float64(radius))
		candidates = m.exclude(candidates, rejected)
		if len(candidates) == 0 {
			radius += 1
			continue
		}

		// availability check
		onlineCandidates := m.availCheck.GetOnlineDrivers(ctx, candidates)
		if len(onlineCandidates) == 0 {
			radius += 1
			continue
		}

		metadatas, err := m.metaFetcher.GetBulkMetada(ctx, onlineCandidates)
		if err != nil {
			return fmt.Errorf("failed to fetch driver metadata: %w", err)
		}

		selectedVehicle := strings.ToUpper(requestDTO.CarType)

		for _, metadata := range metadatas {
			if strings.ToUpper(metadata.VehicleType) == selectedVehicle {
				shortlistDrivers = append(shortlistDrivers, metadata.DriverID)
			}
		}

		if len(shortlistDrivers) > 0 {
			break
		}

		radius += 1
	}

	if len(shortlistDrivers) == 0 {
		payload := domain.NoDriverAvailable{
			TripID:  requestDTO.TripID,
			SagaID:  requestDTO.SagaID,
			Attempt: requestDTO.Attempt,
		}
		evt := ddd.NewEvent(domain.NoDriverAvailableEvent, payload)
		return m.publisher.Publish(ctx, evt)

	}

	replyPayload := &domain.MatchingCandidates{
		SagaID:    requestDTO.SagaID,
		TripID:    requestDTO.TripID,
		DriverIds: shortlistDrivers,
		Attempt:   requestDTO.Attempt,
	}
	matchEvt := ddd.NewEvent(domain.MatchingCandidatesEvent, replyPayload)

	return m.publisher.Publish(ctx, matchEvt)
}
