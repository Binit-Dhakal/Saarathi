package application

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
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

	fullPayload := map[string]any{
		// rider detail
		"rider_id":    riderDetails.ID,
		"rider_name":  riderDetails.Name,
		"rider_phone": riderDetails.PhoneNumber,

		// driver detail
		"driver_id":          driverDetails.ID,
		"driver_name":        driverDetails.Name,
		"driver_phone":       driverDetails.PhoneNumber,
		"vehicle_make":       driverDetails.VehicleMake,
		"vehicle_model":      driverDetails.VehicleModel,
		"license_number":     driverDetails.LicenseNumber,
		"vehicle_number":     driverDetails.VehicleNumber,
		"driver_initial_lat": driverLocation.Lat,
		"driver_initial_lng": driverLocation.Lon,

		// trip detail
		"trip_id":     tripID,
		"pickup_lat":  tripDetail.PickUp.Lat,
		"pickup_lng":  tripDetail.PickUp.Lon,
		"dropoff_lat": tripDetail.DropOff.Lat,
		"dropoff_lng": tripDetail.DropOff.Lon,
		"distance":    tripDetail.Distance,
		"farePrice":   tripDetail.FarePrice,
	}

	return s.saveRepo.SetTripPayload(ctx, tripID, fullPayload, 24*time.Hour)
}
