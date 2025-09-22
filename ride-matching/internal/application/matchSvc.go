package application

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/messagebus"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
)

type MatchingService interface {
	FindDrivers(lon float64, lat float64) []string
}

type matchingService struct {
	publisher am.EventPublisher
	matchRepo domain.RedisRideMatchingRepository
	consumer  messagebus.Consumer
}

func NewMatchingService(publisher am.EventPublisher, matchRepo domain.RedisRideMatchingRepository) MatchingService {
	return &matchingService{
		publisher: publisher,
		matchRepo: matchRepo,
	}
}

// currently our algorithm just find drivers based on geographical location
func (m *matchingService) FindDrivers(lon float64, lat float64) []string {
	return m.matchRepo.FindNearestDriver(lon, lat)
}
