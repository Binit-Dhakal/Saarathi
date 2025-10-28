package application

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/handlers/grpc"
)

type DriverInfoService interface {
	GetDriversMetadata(ctx context.Context, candidates []string) ([]domain.DriverVehicleMetadata, error)
	GetOnlineDrivers(ctx context.Context, candidates []string) []string
}

type driverInfoService struct {
	redisMetaRepo    domain.RedisMetaRepository
	availabilityRepo domain.DriverAvailabilityRepository
	usersClient      *grpc.GRPCClient
}

func NewDriverInfoService(redisMetaRepo domain.RedisMetaRepository, availabilityRepo domain.DriverAvailabilityRepository, client *grpc.GRPCClient) DriverInfoService {
	return &driverInfoService{
		redisMetaRepo:    redisMetaRepo,
		availabilityRepo: availabilityRepo,
		usersClient:      client,
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
		// can be goroutine? need effective solution here
		d.repopulateMetadataCache(missing)

		second_metadata, _ := d.redisMetaRepo.BulkSearchDriverMeta(missing)
		metadatas = append(metadatas, second_metadata...)
	}

	return metadatas, nil
}

func (d *driverInfoService) repopulateMetadataCache(missingCandidates []string) {
	dbMeta, err := d.usersClient.GetBulkDriversMetadata(context.Background(), missingCandidates)
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
