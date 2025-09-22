package dto

import "github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/tripspb"

type TripCreated struct {
	TripID   string
	Distance float64
	Price    int32
	PickUp   *tripspb.Coordinates
	DropOff  *tripspb.Coordinates
	CarType  string
}
