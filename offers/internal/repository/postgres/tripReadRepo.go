package postgres

import (
	"context"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tripReadModelRepository struct {
	pool *pgxpool.Pool
}

func NewTripReadModelRepo(pool *pgxpool.Pool) domain.TripReadModelRepository {
	return &tripReadModelRepository{
		pool: pool,
	}
}

func (t *tripReadModelRepository) SaveTrip(ctx context.Context, payload domain.TripReadModelDTO) error {
	query := `
		INSERT into offers_trip_read_models (trip_id, saga_id,  pickUp, dropOff,distance, price, car_type)
		VALUES($1,$2,point($3,$4),point($5,$6),$7,$8,$9)
	`

	args := []any{payload.TripID, payload.SagaID, payload.PickUp[0], payload.PickUp[1], payload.DropOff[0], payload.DropOff[1], payload.Distance, payload.Price, payload.CarType}
	_, err := t.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return err
}

func (t *tripReadModelRepository) GetTripDetails(ctx context.Context, tripID string) (domain.TripReadModelDTO, error) {
	query := `
		SELECT trip_id, saga_id,pickUp[0],pickUp[1], dropOff[0],dropOff[1], distance, price, car_type from offers_trip_read_models 
		where trip_id=$1
	`

	var (
		result                 domain.TripReadModelDTO
		pickUpLng, pickUpLat   float64
		dropOffLng, dropOffLat float64
	)
	data := t.pool.QueryRow(ctx, query, tripID)

	err := data.Scan(
		&result.TripID, &result.SagaID,
		&pickUpLng, &pickUpLat,
		&dropOffLng, &dropOffLat,
		&result.Distance, &result.Price,
		&result.CarType,
	)
	if err != nil {
		return domain.TripReadModelDTO{}, err
	}
	result.PickUp[0] = pickUpLng
	result.PickUp[1] = pickUpLat
	result.DropOff[0] = dropOffLng
	result.DropOff[1] = dropOffLat

	return result, nil
}
