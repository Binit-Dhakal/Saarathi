package application

import "github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"

type PresenceService interface {
	GetDriverInstance(driverID string) (string, error)
}

type presenceService struct {
	repo domain.PresenceRepository
}

func NewPresenceService(repo domain.PresenceRepository) PresenceService {
	return &presenceService{
		repo: repo,
	}
}

func (p *presenceService) GetDriverInstance(driverID string) (string, error) {
	return p.repo.GetDriverInstanceLocation(driverID)
}
