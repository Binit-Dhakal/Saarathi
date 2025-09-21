package application

//
// import (
// 	"context"
//
// 	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
// 	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
// )
//
// type RideIntegrationService interface {
// 	DriverAccepted(ctx context.Context, input dto.DriverAccepted) error
// }
//
// type rideIntegrationService struct {
// 	tripRepo domain.TripRepository
// }
//
// var _ RideIntegrationService = (*rideIntegrationService)(nil)
//
// func NewRideIntegrationService(tripRepo domain.TripRepository) *rideIntegrationService {
// 	return &rideIntegrationService{
// 		tripRepo: tripRepo,
// 	}
// }
//
// func (r *rideIntegrationService) DriverAccepted(ctx context.Context, input dto.DriverAccepted) error {
// 	return r.tripRepo.UpdateTripDetail(input.TripID, input.DriverID)
// }
