package application

import (
	"context"
	"fmt"

	projectionspb "github.com/Binit-Dhakal/Saarathi/pkg/proto/projections"
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/riderspb"
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

	publicPayload := &riderspb.RiderUpdatePayload{
		TripId:        payload.GetTripId(),
		DriverName:    payload.GetDriverName(),
		VehicleMake:   payload.GetVehicleMake(),
		VehicleModel:  payload.GetVehicleModel(),
		VehicleNumber: payload.GetVehicleNumber(),
		DriverLoc:     payload.DriverLocation,

		Pickup:   payload.Pickup,
		Dropoff:  payload.Dropoff,
		Price:    payload.GetFarePrice(),
		Distance: payload.GetDistance(),
	}
	fmt.Println("Rider payload: ", publicPayload)

	bytesTripPayload, err := proto.Marshal(publicPayload)
	if err != nil {
		return fmt.Errorf("failed to marshal public DTO for trip %s: %w", confirmedDto.TripID, err)
	}

	if s.notifier == nil {
		return fmt.Errorf("notifier is not set for rider service")
	}

	s.notifier.NotifyRider(confirmedDto.TripID, domain.EventSend{
		Event: "TRIP_ACCEPT_PAYLOAD",
		Data:  bytesTripPayload,
	})
	return nil
}
