package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tripRepository struct {
	pool *pgxpool.Pool
}

var _ domain.TripRepository = (*tripRepository)(nil)

func NewTripRepository(pool *pgxpool.Pool) *tripRepository {
	return &tripRepository{
		pool: pool,
	}
}

func (t *tripRepository) SaveRouteDetail(route *domain.Route, riderID string) (string, error) {
	geometryJSON, err := json.Marshal(route.Geometry)
	if err != nil {
		return "", err
	}

	query := `
		INSERT into routes(rider_id, source, destination,distance, duration,geometry) 
		values($1,point($2,$3),point($4,$5),$6,$7,$8)
		returning id
	`

	var routeUUID pgtype.UUID
	var routeID string
	err = t.pool.QueryRow(
		context.Background(),
		query,
		riderID,
		route.Source.Lon, route.Source.Lat,
		route.Destination.Lon, route.Destination.Lat,
		route.Distance,
		route.Duration,
		geometryJSON,
	).Scan(&routeUUID)
	if err != nil {
		return "", err
	}

	routeID = routeUUID.String()

	return routeID, nil
}

func (t *tripRepository) SaveFareDetail(fareModel domain.FareRecord) (string, error) {
	query := `
		INSERT into fares(route_id, car_package, price) 
		VALUES($1,$2,$3)
		returning id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var fareID string
	var fareUUID pgtype.UUID
	args := []any{fareModel.RouteID, fareModel.Fare.Package, fareModel.Fare.TotalPrice}
	err := t.pool.QueryRow(ctx, query, args...).Scan(&fareUUID)
	if err != nil {
		return "", err
	}

	fareID = fareUUID.String()

	return fareID, nil
}

func (t *tripRepository) SaveRideDetail(rideModel domain.RideModel) (string, error) {
	query := `
		INSERT into rides(rider_id,  fare_id, status) 
		values($1,$2,$3)
		returning id
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var rideID string
	var rideUUID pgtype.UUID
	args := []any{rideModel.RiderID, rideModel.FareID, rideModel.Status}
	err := t.pool.QueryRow(ctx, query, args...).Scan(&rideUUID)
	if err != nil {
		return "", err
	}

	rideID = rideUUID.String()

	return rideID, nil
}

func (t *tripRepository) UpdateTripDetail(tripID string, driverID string) error {
	query := `
		UPDATE rides set driver_id=$1 where id=$2 
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.pool.Exec(ctx, query, driverID, tripID)
	if err != nil {
		return err
	}

	return nil
}
