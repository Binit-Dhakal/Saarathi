package waiter

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Runner struct {
	Timeout time.Duration
}

func NewRunner(timeout time.Duration) *Runner {
	return &Runner{
		Timeout: timeout,
	}
}

func (r *Runner) RunHTTPServer(ctx context.Context, s *http.Server) error {
	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Printf("HTTP Server started; listening at %s\n", s.Addr)
		defer fmt.Println("HTTP Server shutdown")
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("http server failed: %w", err)
		}
		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("HTTP Server initiating shutdown")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), r.Timeout)
		defer cancel()

		if err := s.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("http server forced shutdown: %w", err)
		}
		return nil
	})

	return group.Wait()
}

func (r *Runner) RunGRPCServer(ctx context.Context, s *grpc.Server, addr string) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on gRPC address %s: %w", addr, err)
	}

	group, gCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		fmt.Printf("gRPC server started; listening at %s\n", addr)
		defer fmt.Println("gRPC server shutdown")

		if err := s.Serve(listener); err != nil && status.Code(err) != codes.Unavailable {
			return fmt.Errorf("grpc server failed: %w", err)
		}

		return nil
	})

	group.Go(func() error {
		<-gCtx.Done()
		fmt.Println("gRPC server initiating shutdown")

		stopped := make(chan struct{})
		go func() {
			s.GracefulStop()
			close(stopped)
		}()

		timeout := time.NewTimer(r.Timeout)
		select {
		case <-timeout.C:
			s.Stop()
			return fmt.Errorf("gRPC server failed to stop gracefully within timeout")
		case <-stopped:
			return nil
		}
	})

	return group.Wait()
}
