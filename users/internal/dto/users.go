package dto

type RiderRegistrationDTO struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phoneNumber"`
}

type DriverRegistrationDTO struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PhoneNumber   string `json:"phoneNumber"`
	LicenseNumber string `json:"licenseNumber"`
	VehicleNumber string `json:"vehicleNumber"`
	VehicleModel  string `json:"vehicleModel"`
	VehicleMake   string `json:"vehicleMake"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
