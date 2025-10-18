package application

import (
	"context"
	"encoding/json"
	"fmt"

	projectionspb "github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/projections"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/domain"
	"google.golang.org/protobuf/proto"
)

type RiderUpdateService interface {
	AddNotifier(notifier domain.Notifier)
	SendTripCompleteDetail(ctx context.Context, confirmedDto *domain.TripConfirmedDTO) error
}

type riderUpdateService struct {
	notifier domain.Notifier
	repo     domain.TripPayloadRepository
}

var _ RiderUpdateService = (*riderUpdateService)(nil)

func NewRiderUpdateService(repo domain.TripPayloadRepository) RiderUpdateService {
	return &riderUpdateService{
		repo: repo,
	}
}

func (s *riderUpdateService) AddNotifier(notifier domain.Notifier) {
	s.notifier = notifier
}

func (s *riderUpdateService) SendTripCompleteDetail(ctx context.Context, confirmedDto *domain.TripConfirmedDTO) error {
	bytesPayload, err := s.repo.GetTripFullPayload(ctx, confirmedDto.TripID)
	if err != nil {
		return err
	}

	payload := &projectionspb.TripProjectionV1{}

	if err := proto.Unmarshal(bytesPayload, payload); err != nil {
		return fmt.Errorf("failed to unmarshal trip payload into Protobuf DTO for trip %s: %w", confirmedDto.TripID, err)
	}

	publicPayload := &domain.RiderUpdateDTO{
		TripID:        payload.GetTripId(),
		DriverName:    payload.GetDriverName(),
		VehicleMake:   payload.GetVehicleMake(),
		VehicleModel:  payload.GetVehicleModel(),
		VehicleNumber: payload.GetVehicleNumber(),
		DriverLat:     payload.DriverLocation.GetLat(),
		DriverLng:     payload.DriverLocation.GetLng(),

		PickupLat:  payload.Pickup.GetLat(),
		PickupLng:  payload.Pickup.GetLng(),
		DropoffLat: payload.Dropoff.GetLat(),
		DropoffLng: payload.Dropoff.GetLng(),
		FarePrice:  payload.GetFarePrice(),
		Distance:   payload.GetDistance(),
	}

	jsonBytes, err := json.Marshal(publicPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal public DTO for trip %s: %w", confirmedDto.TripID, err)
	}

	if s.notifier == nil {
		fmt.Errorf("notifier is not set for rider service")
	}

	s.notifier.NotifyRider(confirmedDto.TripID, jsonBytes)
	return nil
}
