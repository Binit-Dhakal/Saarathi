package dto

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/proto/common"
)

type TripCreated struct {
	SagaID           string
	TripID           string
	PickUp           *common.Coordinates
	DropOff          *common.Coordinates
	CarType          string
	SearchRadius     int32
	Attempt          int32
	FirstAttemptUnix int64
}
