package domain

import (
	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
)

type Driver struct {
	ddd.Aggregate
	Available bool
	Location  Location
	Offers    map[string]*Offer
}
