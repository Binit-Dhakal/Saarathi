package application

import "github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"

type PresenceService interface {
	SaveWSDetail(driverID string) error
	DeleteWSDetail(driverID string) error
}

type presenceSvc struct {
	wsRepo domain.WSRepo
}

func NewPresenceService(wsRepo domain.WSRepo) PresenceService {
	return &presenceSvc{
		wsRepo: wsRepo,
	}
}

func (p *presenceSvc) SaveWSDetail(driverID string) error {
	return p.wsRepo.SaveWSDetail(driverID)
}

func (p *presenceSvc) DeleteWSDetail(driverID string) error {
	return p.wsRepo.DeleteWSDetail(driverID)
}
