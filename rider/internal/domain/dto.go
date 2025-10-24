package domain

type TripConfirmedDTO struct {
	TripID   string
	DriverID string
	RiderID  string
}

type RiderUpdateDTO struct {
	TripID string `json:"trip_id"`

	// Driver Info (Minimal)
	DriverName    string  `json:"driver_name"`
	VehicleMake   string  `json:"vehicle_make"`
	VehicleModel  string  `json:"vehicle_model"`
	VehicleNumber string  `json:"vehicle_number"`
	DriverLat     float64 `json:"driver_lat"`
	DriverLng     float64 `json:"driver_lng"`

	// Trip Info (Only what the rider sees)
	PickupLat  float64 `json:"pickup_lat"`
	PickupLng  float64 `json:"pickup_lng"`
	DropoffLat float64 `json:"dropoff_lat"`
	DropoffLng float64 `json:"dropoff_lng"`
	FarePrice  int32   `json:"fare_price"`
	Distance   float64 `json:"distance"`
}

type EventSend struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}
