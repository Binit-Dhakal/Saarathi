package dto

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/contracts/proto/common"
)

type TripCreated struct {
	TripID   string
	Distance float64
	Price    int32
	PickUp   *common.Coordinates
	DropOff  *common.Coordinates
	CarType  string
}
