package application

import "context"

type ProjectionService interface {
	ProjectTripDetails(ctx context.Context, tripID, driverID, riderID string) error
}

type projectionService struct {
}

func NewProjectionService() ProjectionService {
	return &projectionService{}
}

func (p *projectionService) ProjectTripDetails(ctx context.Context, tripID, driverID, riderID string) error {
	return nil
}
