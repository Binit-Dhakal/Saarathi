package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/users/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUniqueViolation = errors.New("Unique Constraint Violated")
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

func (u *UserRepo) CreateUser(tx pgx.Tx, user *domain.User) (string, error) {
	query := `
		INSERT into users(name, email, country, phone_number, password) 
		VALUES ($1, $2,$3,$4,$5)
		RETURNING id
	`
	hashedPassword, err := domain.Hash(user.Password)
	if err != nil {
		return "", err
	}

	args := []any{user.Name, user.Email, user.Country, user.PhoneNumber, hashedPassword}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userUUID pgtype.UUID
	err = tx.QueryRow(ctx, query, args...).Scan(&userUUID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505":
				return "", ErrUniqueViolation
			default:
				return "", fmt.Errorf("PostgreSQL error: %v (Code: %s)\n", pgErr.Message, pgErr.Code)
			}
		}
		return "", err
	}

	return userUUID.String(), nil
}

func (u *UserRepo) AddUserToRole(tx pgx.Tx, userID string, role int) error {
	query := `
		INSERT into user_roles(user_id, role_id) VALUES($1,$2)
	`
	args := []any{userID, role}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.Exec(ctx, query, args...)
	return err
}

func (u *UserRepo) CreateRiderProfile(tx pgx.Tx, profile *domain.RiderProfile) error {
	// TODO: deal with payment in the rider profile
	query := `
		INSERT into rider_profiles(user_id) VALUES ($1)
	`
	args := []any{profile.UserID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return ErrUniqueViolation
			default:
				return fmt.Errorf("PostgreSQL error: %v (Code: %s)\n", pgErr.Message, pgErr.Code)
			}
		}
		return err
	}

	return nil
}

func (u *UserRepo) CreateDriverProfile(tx pgx.Tx, profile *domain.DriverProfile) error {
	query := `
		INSERT into driver_profiles(user_id, license_number, vehicle_number, vehicle_model, vehicle_make) 
		VALUES ($1,$2,$3,$4,$5)
	`
	args := []any{profile.UserID, profile.LicenseNumber, profile.VehicleNumber, profile.VehicleModel, profile.VehicleMake}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := tx.Exec(ctx, query, args...)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch {
			case pgErr.Code == "23505":
				return ErrUniqueViolation
			default:
				return fmt.Errorf("PostgreSQL error: %v (Code: %s)\n", pgErr.Message, pgErr.Code)
			}
		}
		return err
	}

	return nil
}

func (u *UserRepo) GetUserByEmail(tx pgx.Tx, email string) (*domain.User, error) {
	query := `SELECT id, email, password from users where email=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user domain.User
	var userID pgtype.UUID
	err := tx.QueryRow(ctx, query, email).Scan(&userID, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	user.ID = userID.String()
	return &user, nil
}

func (u *UserRepo) GetForToken(refreshToken string) (*domain.User, error) {
	query := `SELECT user_id from tokens where refresh_token=$1`
	var user domain.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userID pgtype.UUID
	err := u.pool.QueryRow(ctx, query, refreshToken).Scan(&userID)
	if err != nil {
		return nil, err
	}

	user.ID = userID.String()

	return &user, nil
}
