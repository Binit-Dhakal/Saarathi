package setup

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetupRpc() *grpc.Server {
	server := grpc.NewServer()
	reflection.Register(server)

	return server
}
