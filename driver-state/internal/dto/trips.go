package dto

type TripAssignedDTO struct {
	TripID   string
	DriverID string
	RiderID  string
}

type DriverUpdateDTO struct {
	TripID string `json:"trip_id"`

	// Rider Info
	RiderName   string `json:"rider_name"`
	RiderNumber string `json:"rider_phone"`

	// Trip Info (Only what the rider sees)
	PickupLat  float64 `json:"pickup_lat"`
	PickupLng  float64 `json:"pickup_lng"`
	DropoffLat float64 `json:"dropoff_lat"`
	DropoffLng float64 `json:"dropoff_lng"`
	FarePrice  int32   `json:"fare_price"`
	Distance   float64 `json:"distance"`
}
