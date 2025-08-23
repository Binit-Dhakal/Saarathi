package application

import (
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
)

type LocationService interface {
	UpsertDriverLocation(loc *dto.DriverLocationMessage, driverID string) error
}

type locationService struct {
	locationRepo domain.LocationRepo
}

func NewLocationService(repo domain.LocationRepo) *locationService {
	return &locationService{
		locationRepo: repo,
	}
}

func (l *locationService) UpsertDriverLocation(loc *dto.DriverLocationMessage, driverID string) error {
	locInfo := &domain.DriverLocation{
		DriverID:    driverID,
		Longitude:   loc.Longitude,
		Latitude:    loc.Latitude,
		VehicleType: loc.CarPackage,
	}

	err := l.locationRepo.SaveActiveGeoLocation(locInfo)
	if err != nil {
		return err
	}

	return nil
}
