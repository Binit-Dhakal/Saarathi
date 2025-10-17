package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/userspb"
	"github.com/Binit-Dhakal/Saarathi/trips/internals/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	addr   string
	conn   *grpc.ClientConn
	client userspb.UsersServiceClient
}

func NewGRPCClient(addr string) (*GRPCClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Users service at %s: %w", addr, err)
	}

	return &GRPCClient{
		addr:   addr,
		conn:   conn,
		client: userspb.NewUsersServiceClient(conn),
	}, nil
}

func (c *GRPCClient) Close() error {
	return c.conn.Close()
}

func (c *GRPCClient) GetRiderDetails(ctx context.Context, riderID string) (*domain.RiderDetail, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.GetRiderDetails(rpcCtx, &userspb.GetRiderDetailsRequest{Id: riderID})
	if err != nil {
		return nil, fmt.Errorf("grpc error fetching rider %s: %w", riderID, err)
	}

	return &domain.RiderDetail{
		ID:          resp.Id,
		Name:        resp.Name,
		PhoneNumber: resp.PhoneNumber,
	}, nil
}

func (c *GRPCClient) GetDriverDetails(ctx context.Context, driverID string) (*domain.DriverDetail, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.GetDriverDetails(rpcCtx, &userspb.GetDriverDetailsRequest{Id: driverID})
	if err != nil {
		return nil, fmt.Errorf("grpc error fetching driver %s: %w", driverID, err)
	}

	return &domain.DriverDetail{
		ID:            resp.Id,
		Name:          resp.Name,
		PhoneNumber:   resp.PhoneNumber,
		LicenseNumber: resp.LicenseNumber,
		VehicleMake:   resp.VehicleMake,
		VehicleModel:  resp.VehicleModel,
		VehicleNumber: resp.VehiclePlate,
	}, nil
}
