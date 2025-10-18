package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/common"
	projectionspb "github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/projections"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"google.golang.org/protobuf/proto"
)

type ProjectionService interface {
	ProjectTripDetails(ctx context.Context, tripID, driverID, riderID string) error
}

type projectionService struct {
	usersClient     domain.UsersClient
	presenceGateway domain.PresenceGatewayRepository
	saveRepo        domain.TripProjectionRepository
	tripReadRepo    domain.TripReadRepository
}

func NewProjectionService(usersClient domain.UsersClient, presenceGateway domain.PresenceGatewayRepository, saveRepo domain.TripProjectionRepository, tripReadRepo domain.TripReadRepository) ProjectionService {
	return &projectionService{
		usersClient:     usersClient,
		presenceGateway: presenceGateway,
		saveRepo:        saveRepo,
		tripReadRepo:    tripReadRepo,
	}
}

func (s *projectionService) ProjectTripDetails(ctx context.Context, tripID, driverID, riderID string) error {
	driverLocation, err := s.presenceGateway.GetDriverLocation(ctx, driverID)
	if err != nil {
		return err
	}

	tripDetail, err := s.tripReadRepo.GetTripProjectionDetail(ctx, tripID)
	if err != nil {
		return err
	}

	driverDetails, err := s.usersClient.GetDriverDetails(ctx, driverID)
	if err != nil {
		return fmt.Errorf("failed to fetch driver details: %w", err)
	}

	riderDetails, err := s.usersClient.GetRiderDetails(ctx, riderID)
	if err != nil {
		return fmt.Errorf("failed to fetch rider details: %w", err)
	}

	fullPayload := projectionspb.TripProjectionV1{
		// rider detail
		RiderId:    riderDetails.ID,
		RiderName:  riderDetails.Name,
		RiderPhone: riderDetails.PhoneNumber,

		// driver detail
		DriverId:       driverDetails.ID,
		DriverName:     driverDetails.Name,
		DriverPhone:    driverDetails.PhoneNumber,
		VehicleMake:    driverDetails.VehicleMake,
		VehicleModel:   driverDetails.VehicleModel,
		LicenseNumber:  driverDetails.LicenseNumber,
		VehicleNumber:  driverDetails.VehicleNumber,
		DriverLocation: &common.Coordinates{Lng: driverLocation.Lon, Lat: driverLocation.Lat},

		// trip detail
		TripId:    tripID,
		Pickup:    &common.Coordinates{Lng: tripDetail.PickUp.Lon, Lat: tripDetail.PickUp.Lat},
		Dropoff:   &common.Coordinates{Lng: tripDetail.DropOff.Lon, Lat: tripDetail.DropOff.Lat},
		Distance:  tripDetail.Distance,
		FarePrice: int32(tripDetail.FarePrice),
	}

	bytePayload, err := proto.Marshal(&fullPayload)
	if err != nil {
		return err
	}

	return s.saveRepo.SetTripPayload(ctx, tripID, bytePayload, 24*time.Hour)
}
