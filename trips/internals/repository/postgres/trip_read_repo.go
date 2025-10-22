package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type tripReadRepository struct {
	pool *pgxpool.Pool
}

func NewTripReadRepository(pool *pgxpool.Pool) *tripReadRepository {
	return &tripReadRepository{
		pool: pool,
	}
}

func (t *tripReadRepository) GetTripProjectionDetail(ctx context.Context, tripID string) (domain.TripProjectionDetail, error) {
	query := `
		SELECT
			r.id,
			f.price,
			rt.distance,
			rt.source[0],
			rt.source[1],
			rt.destination[0],
			rt.destination[1]
		FROM rides r 
		JOIN fares f on r.fare_id = f.id
		JOIN routes rt on f.route_id = rt.id 
		WHERE r.id = $1
	`
	var detail domain.TripProjectionDetail

	row := t.pool.QueryRow(ctx, query, tripID)
	err := row.Scan(
		&detail.TripID,
		&detail.FarePrice,
		&detail.Distance,
		&detail.PickUp.Lon,
		&detail.PickUp.Lat,
		&detail.DropOff.Lon,
		&detail.DropOff.Lat,
	)
	if err == sql.ErrNoRows {
		return detail, fmt.Errorf("trip projection detail not found for trip ID: %s", tripID)
	}

	if err != nil {
		return detail, fmt.Errorf("error querying trip projection detail: %w", err)
	}

	return detail, nil
}
