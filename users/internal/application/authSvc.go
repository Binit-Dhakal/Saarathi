package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/users/internal/domain"
	"github.com/Binit-Dhakal/Saarathi/users/internal/dto"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrWrongCredentials = errors.New("Wrong Credentials")
)

type AuthService interface {
	RegisterRider(rider *dto.RiderRegistrationDTO) (userID string, err error)
	RegisterDriver(driver *dto.DriverRegistrationDTO) (userID string, err error)
	CreateAuthenticationToken(creds *dto.LoginRequestDTO) (userID string, err error)
	// RefreshToken(refreshToken string) (*Token, error)
}

type AuthServiceImpl struct {
	pool     *pgxpool.Pool
	userRepo domain.UserRepo
}

func NewAuthService(pool *pgxpool.Pool, userRepo domain.UserRepo) *AuthServiceImpl {
	return &AuthServiceImpl{
		pool:     pool,
		userRepo: userRepo,
	}
}

func (a *AuthServiceImpl) RegisterRider(rider *dto.RiderRegistrationDTO) (userID string, err error) {
	ctx := context.Background()
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	dtoUser := &domain.User{
		Name:        rider.Name,
		Email:       rider.Email,
		Password:    rider.Password,
		PhoneNumber: rider.PhoneNumber,
	}
	userID, err = a.userRepo.CreateUser(tx, dtoUser)
	if err != nil {
		return "", err
	}

	riderProfile := &domain.RiderProfile{
		UserID: userID,
	}
	err = a.userRepo.CreateRiderProfile(tx, riderProfile)
	if err != nil {
		return "", err
	}

	err = a.userRepo.AddUserToRole(tx, userID, domain.RoleRider)
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (a *AuthServiceImpl) RegisterDriver(driver *dto.DriverRegistrationDTO) (userID string, err error) {
	ctx := context.Background()
	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	dtoUser := &domain.User{
		Name:        driver.Name,
		Email:       driver.Email,
		Password:    driver.Password,
		PhoneNumber: driver.PhoneNumber,
	}
	userID, err = a.userRepo.CreateUser(tx, dtoUser)
	if err != nil {
		return "", err
	}

	driverProfile := &domain.DriverProfile{
		UserID:        userID,
		LicenseNumber: driver.LicenseNumber,
		VehicleNumber: driver.VehicleNumber,
		VehicleMake:   driver.VehicleMake,
		VehicleModel:  driver.VehicleModel,
	}
	err = a.userRepo.CreateDriverProfile(tx, driverProfile)
	if err != nil {
		return "", err
	}

	err = a.userRepo.AddUserToRole(tx, userID, domain.RoleDriver)
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}

	return userID, nil
}

func (a *AuthServiceImpl) CreateAuthenticationToken(creds *dto.LoginRequestDTO) (userID string, err error) {
	tx, err := a.pool.Begin(context.Background())
	if err != nil {
		return "", err
	}
	user, err := a.userRepo.GetUserByEmail(tx, creds.Email)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	ok, err := domain.Matches(creds.Password, user.Password)
	if err != nil {
		fmt.Println(err)
	}
	if !ok {
		return "", ErrWrongCredentials
	}

	return user.ID, nil
}
