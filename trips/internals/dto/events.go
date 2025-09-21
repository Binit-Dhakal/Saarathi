package dto

type AcceptDriver struct {
	DriverID string
	TripID   string
}

type TripConfirmed struct {
	TripID        string `json:"tripID"`
	DriverID      string `json:"driverID"`
	DriverName    string `json:"driverName"`
	VehicleNumber string `json:"vehicleNumber"`
	ContactNumber string `json:"contactNumber"`
	Status        string `json:"status"`
}
