package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/userspb"
	"github.com/Binit-Dhakal/Saarathi/users/internal/application"
	"google.golang.org/grpc"
)

var (
	ErrRoleMismatch = errors.New("Role not matched")
)

type server struct {
	svc application.UserDetailService
	userspb.UnimplementedUsersServiceServer
}

var _ userspb.UsersServiceServer = (*server)(nil)

func RegisterServer(svc application.UserDetailService, registrar grpc.ServiceRegistrar) error {
	userspb.RegisterUsersServiceServer(registrar, server{svc: svc})
	return nil
}

func (s server) GetRiderDetails(ctx context.Context, req *userspb.GetRiderDetailsRequest) (*userspb.GetRiderDetailsResponse, error) {
	user, err := s.svc.GetRiderDetails(ctx, req.Id)
	if err != nil {
		// TODO: need to handle not found and server error
		return nil, fmt.Errorf("failed to retrieve error: %w", err)
	}

	if user.Role != "rider" {
		return nil, ErrRoleMismatch
	}

	return &userspb.GetRiderDetailsResponse{
		Id:          user.ID,
		Name:        user.Name,
		PhoneNumber: user.PhoneNumber,
	}, nil
}

func (s server) GetDriverDetails(ctx context.Context, req *userspb.GetDriverDetailsRequest) (*userspb.GetDriverDetailsResponse, error) {
	user, err := s.svc.GetDriverDetails(ctx, req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve error: %w", err)
	}

	if user.Role != "driver" {
		return nil, ErrRoleMismatch
	}

	return &userspb.GetDriverDetailsResponse{
		Id:            user.ID,
		Name:          user.Name,
		PhoneNumber:   user.PhoneNumber,
		LicenseNumber: user.LicenseNumber,
		VehicleMake:   user.VehicleMake,
		VehicleModel:  user.VehicleModel,
		VehiclePlate:  user.VehicleNumber,
	}, nil
}
