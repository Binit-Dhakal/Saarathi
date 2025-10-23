package dto

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/common"
)

type TripCreated struct {
	SagaID       string
	TripID       string
	PickUp       *common.Coordinates
	DropOff      *common.Coordinates
	CarType      string
	SearchRadius float64
}
