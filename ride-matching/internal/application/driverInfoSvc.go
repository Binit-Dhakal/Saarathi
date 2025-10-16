package application

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
)

type DriverInfoService interface {
	GetDriversMetadata(ctx context.Context, candidates []string) ([]domain.DriverVehicleMetadata, error)
	GetOnlineDrivers(ctx context.Context, candidates []string) []string
}

type driverInfoService struct {
	redisMetaRepo    domain.RedisMetaRepository
	pgMetaRepo       domain.PGMetaRepository
	availabilityRepo domain.DriverAvailabilityRepository
}

func NewDriverInfoService(redisMetaRepo domain.RedisMetaRepository, pgMetaRepo domain.PGMetaRepository, availabilityRepo domain.DriverAvailabilityRepository) DriverInfoService {
	return &driverInfoService{
		redisMetaRepo:    redisMetaRepo,
		pgMetaRepo:       pgMetaRepo,
		availabilityRepo: availabilityRepo,
	}
}

func (d *driverInfoService) GetDriversMetadata(ctx context.Context, candidates []string) ([]domain.DriverVehicleMetadata, error) {
	metadatas, err := d.redisMetaRepo.BulkSearchDriverMeta(candidates)
	if err != nil {
		return nil, err
	}

	var missing []string
	for _, d := range metadatas {
		if d.VehicleType == "" {
			missing = append(missing, d.DriverID)
		}
	}
	if len(missing) > 0 {
		go d.repopulateMetadataCache(missing)
	}

	return metadatas, nil
}

func (d *driverInfoService) repopulateMetadataCache(missingCandidates []string) {
	dbMeta, err := d.pgMetaRepo.BulkSearchMeta(missingCandidates)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = d.redisMetaRepo.BulkInsertDriverMeta(dbMeta)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (d *driverInfoService) GetOnlineDrivers(ctx context.Context, candidates []string) []string {
	valid, expired := d.availabilityRepo.BulkCheckDriversOnline(candidates)

	d.availabilityRepo.DeleteUnavailableDrivers(expired)

	return valid
}
