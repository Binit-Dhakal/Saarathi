package application

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/common"
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
	"github.com/google/uuid"
)

type RideService interface {
	EstimateFare(route *domain.Route) ([]domain.Fare, string, error)
	FareAcceptByRider(req *dto.FareConfirmRequest, userID string) (string, error)
}

type rideService struct {
	fareRepo   domain.FareRepository
	tripRepo   domain.TripRepository
	tripStream am.EventPublisher
	sagaStream am.EventPublisher
}

func NewRideService(fareRepo domain.FareRepository, tripRepo domain.TripRepository, tripStream am.EventPublisher, sagaStream am.EventPublisher) *rideService {
	return &rideService{
		fareRepo:   fareRepo,
		tripRepo:   tripRepo,
		tripStream: tripStream,
		sagaStream: sagaStream,
	}
}

func (f *rideService) EstimateFare(route *domain.Route) ([]domain.Fare, string, error) {
	distanceInKm := (route.Distance) / 1000
	etaInMinutes := (route.Distance) / 60

	if distanceInKm == 0 {
		distanceInKm = 1
	}

	if etaInMinutes == 0 {
		etaInMinutes = 1
	}

	var fares []domain.Fare
	for _, fare := range domain.CarRegistry {
		totalPrice := fare.BaseFare + (fare.PerKmRate * int(distanceInKm)) + (fare.PerMinuteRate * int(etaInMinutes))
		fares = append(
			fares,
			domain.Fare{
				Package:    fare.Name,
				TotalPrice: totalPrice,
			},
		)
	}

	// save the fare detail to redis and get fareID to later retrieve
	ephemeralFare := domain.FareQuote{
		Route: *route,
		Fares: fares,
	}

	fareID, err := f.fareRepo.CreateEphemeralFareEntry(&ephemeralFare)
	if err != nil {
		return nil, "", err
	}

	return fares, fareID, nil
}

// 1. Get the ephermal fare back from redis
// 2. store the route for given fare
// 3. Store the fare for given carpackage in fares table
// 4. persist ride (rider + fare + driver(eventually))
func (f *rideService) FareAcceptByRider(req *dto.FareConfirmRequest, userID string) (string, error) {
	ephermalFare, err := f.fareRepo.GetEphemeralFareEntry(req.FareID)
	if err != nil {
		return "", err
	}

	routeID, err := f.tripRepo.SaveRouteDetail(&ephermalFare.Route, userID)
	if err != nil {
		return "", err
	}

	var fareDetail domain.Fare
	for _, fare := range ephermalFare.Fares {
		if fare.Package == req.CarPackage {
			fareDetail = fare
			break
		}
	}

	fareRecord := domain.FareRecord{
		Fare:    fareDetail,
		RouteID: routeID,
	}

	fareID, err := f.tripRepo.SaveFareDetail(fareRecord)
	if err != nil {
		return "", err
	}

	rideModel := domain.TripModel{
		RiderID: userID,
		FareID:  fareID,
		Status:  domain.TripStatusPending,
	}

	rideId, err := f.tripRepo.SaveRideDetail(rideModel)
	if err != nil {
		return "", err
	}

	createdEvent := tripspb.TripRequested{
		SagaId:   uuid.NewString(),
		TripId:   rideId,
		Distance: ephermalFare.Route.Distance,
		Price:    int32(fareDetail.TotalPrice),
		PickUp:   &common.Coordinates{Lng: ephermalFare.Route.Source.Lon, Lat: ephermalFare.Route.Source.Lat},
		DropOff:  &common.Coordinates{Lng: ephermalFare.Route.Destination.Lon, Lat: ephermalFare.Route.Destination.Lat},
		CarType:  string(fareDetail.Package),
	}

	event := ddd.NewEvent(tripspb.TripRequestedEvent, &createdEvent)
	err = f.tripStream.Publish(
		context.Background(),
		tripspb.TripRequestedEvent,
		event,
	)
	if err != nil {
		return "", err
	}

	err = f.sagaStream.Publish(
		context.Background(),
		tripspb.TripRequestedEvent,
		event,
	)
	if err != nil {
		return "", err
	}

	return rideId, nil
}
