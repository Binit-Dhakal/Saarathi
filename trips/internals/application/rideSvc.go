package application

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/pkg/am"
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/dto"
	tripsv1 "github.com/Binit-Dhakal/Saarathi/trips/tripspb/proto/trip"
)

type RideService interface {
	EstimateFare(route *domain.Route) ([]domain.Fare, string, error)
	FareAcceptByRider(req *dto.FareConfirmRequest, userID string) (string, error)
}

type rideService struct {
	fareRepo domain.FareRepository
	tripRepo domain.TripRepository
	bus      am.EventPublisher
}

func NewRideService(fareRepo domain.FareRepository, tripRepo domain.TripRepository, bus am.EventPublisher) *rideService {
	return &rideService{
		fareRepo: fareRepo,
		tripRepo: tripRepo,
		bus:      bus,
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

	rideModel := domain.RideModel{
		RiderID: userID,
		FareID:  fareID,
		Status:  domain.RideStatusPending,
	}

	rideId, err := f.tripRepo.SaveRideDetail(rideModel)
	if err != nil {
		return "", err
	}

	createdEvent := tripsv1.TripCreated{
		TripId:   rideId,
		Distance: ephermalFare.Route.Distance,
		Price:    int32(fareDetail.TotalPrice),
		PickUp:   &tripsv1.Coordinates{Lng: ephermalFare.Route.Source.Lon, Lat: ephermalFare.Route.Source.Lat},
		DropOff:  &tripsv1.Coordinates{Lng: ephermalFare.Route.Destination.Lon, Lat: ephermalFare.Route.Destination.Lat},
		CarType:  string(fareDetail.Package),
	}

	event := ddd.NewEvent("trips.created", &createdEvent)
	err = f.bus.Publish(
		context.Background(),
		event.EventName(),
		event,
	)
	if err != nil {
		return "", err
	}
	return rideId, nil
}
