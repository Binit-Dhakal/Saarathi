package postgres

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type metaRepo struct {
	pool *pgxpool.Pool
}

func NewPGMetaRepo(pool *pgxpool.Pool) domain.PGMetaRepository {
	return &metaRepo{
		pool: pool,
	}
}

func (m *metaRepo) BulkSearchMeta(driverIDs []string) ([]domain.DriverVehicleMetadata, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		select u.id, dp.vehicle_model 
		from users u
		join driver_profiles dp ON dp.user_id = u.id
		WHERE u.id = ANY($1)
	`

	rows, err := m.pool.Query(ctx, query, driverIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByPos[domain.DriverVehicleMetadata])
}

