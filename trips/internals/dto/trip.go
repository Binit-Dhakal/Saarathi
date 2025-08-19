package dto

import "github.com/Binit-Dhakal/Saarathi/trips/internals/domain"

type FareEstimateRequest struct {
	PickUpLocation  []float64 `json:"pickUpLocation"`
	DropOffLocation []float64 `json:"dropOffLocation"`
}

type FareEstimateResponse struct {
	FareID   string
	Fares    []domain.Fare
	Geometry domain.Geometry
}

type FareConfirmRequest struct {
	FareID     string            `json:"fareID"`
	CarPackage domain.CarPackage `json:"carPackage"`
}

type FareConfirmResponse struct {
	RideID string `json:"rideID"`
}
