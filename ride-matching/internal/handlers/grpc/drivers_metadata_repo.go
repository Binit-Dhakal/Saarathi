package grpc

import (
	"context"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/proto/userspb"
	"github.com/Binit-Dhakal/Saarathi/ride-matching/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	addr   string
	conn   *grpc.ClientConn
	client userspb.UsersServiceClient
	userspb.UnimplementedUsersServiceServer
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

func (c *GRPCClient) GetBulkDriversMetadata(ctx context.Context, driver_ids []string) ([]domain.DriverVehicleMetadata, error) {
	rpcCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	resp, err := c.client.GetDriverMetadataBatch(rpcCtx, &userspb.GetDriverMetadataBatchRequest{Ids: driver_ids})
	if err != nil {
		return nil, fmt.Errorf("grpc error fetching driver metadata %v:%w", driver_ids, err)
	}

	result := []domain.DriverVehicleMetadata{}
	for _, r := range resp.Drivers {
		result = append(result, domain.DriverVehicleMetadata{
			DriverID:    r.Id,
			VehicleType: r.VehicleModel,
		})

	}
	return result, nil
}
