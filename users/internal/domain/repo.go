package domain

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type UserRepo interface {
	CreateUser(tx pgx.Tx, user *User) (string, error)
	AddUserToRole(tx pgx.Tx, userID string, role int) error
	CreateRiderProfile(tx pgx.Tx, profile *RiderProfile) error
	CreateDriverProfile(tx pgx.Tx, profile *DriverProfile) error
	GetUserByEmail(tx pgx.Tx, email string) (*User, error)

	GetRiderByID(ctx context.Context, id string) (*UserDetail, error)
	GetDriverByID(ctx context.Context, id string) (*DriverDetail, error)

	BulkSearchMeta(ctx context.Context, driverIDs []string) ([]DriverVehicleMetadata, error)
}

type TokenRepo interface {
	CreateToken(token *Token) error
	FindByRefreshToken(refreshToken string) (*Token, error)
	RevokeRefreshToken(refreshToken string) error
}
