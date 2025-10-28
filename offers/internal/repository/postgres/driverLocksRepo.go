package postgres

import (
	"context"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/offers/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type driverLocksRepository struct {
	pool *pgxpool.Pool
}

func NewDriverLocksRepository(pool *pgxpool.Pool) domain.DriverLockRepository {
	return &driverLocksRepository{
		pool: pool,
	}
}

func (d *driverLocksRepository) TryAcquireLock(ctx context.Context, driverID string, tripID string) (bool, error) {
	query := `
		INSERT INTO driver_locks (driver_id, status, trip_id, expired_at)
		VALUES ($1, 'OFFERED', $2, NOW() + interval '15 seconds')
		ON CONFLICT (driver_id) DO UPDATE 
		SET status = 'OFFERED',
			trip_id = EXCLUDED.trip_id,
			expired_at = NOW()+interval '15 seconds'
		WHERE driver_locks.status = 'AVAILABLE'
			OR (driver_locks.status = 'OFFERED' AND driver_locks.expired_at < NOW())
	`

	cmd, err := d.pool.Exec(ctx, query, driverID, tripID)
	if err != nil {
		return false, fmt.Errorf("Upsert failed: %w", err)
	}

	return cmd.RowsAffected() > 0, nil
}

func (d *driverLocksRepository) ReleaseLock(ctx context.Context, driverID string, tripID string) error {
	query := `
		UPDATE driver_locks 
		SET status='AVAILABLE', trip_id=NULL, expired_at = NULL
		WHERE driver_id=$1 and trip_id = $2
	`
	cmd, err := d.pool.Exec(ctx, query, driverID, tripID)
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		// silent ignore - not owned by trip - but maybe we should do something about this?
	}
	return nil
}

func (d *driverLocksRepository) AcceptTrip(ctx context.Context, driverID, tripID string) error {
	query := `
		UPDATE driver_locks 
		SET status='ACCEPTED'
		WHERE driver_id = $1 and trip_id =$2
	`

	_, err := d.pool.Exec(ctx, query, driverID, tripID)
	if err != nil {
		return fmt.Errorf("failed to mark driver as accepted: %w", err)
	}

	return nil
}

func (d *driverLocksRepository) MarkAvailable(ctx context.Context, driverID string) error {
	query := `
		UPDATE driver_locks
		SET status='AVAILABLE', trip_id = NULL, expired_at = NULL
		WHERE driver_id = $1
	`
	_, err := d.pool.Exec(ctx, query, driverID)
	if err != nil {
		return fmt.Errorf("failed to mark driver as available: %w", err)
	}

	return nil
}
