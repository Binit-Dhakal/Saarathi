package domain

type RiderDetail struct {
	ID          string
	Name        string
	PhoneNumber string
}

type DriverDetail struct {
	ID            string
	Name          string
	PhoneNumber   string
	LicenseNumber string
	VehicleMake   string
	VehicleModel  string
	VehicleNumber string
}
