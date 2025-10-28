package domain

import "time"

type User struct {
	ID          string
	Name        string
	Email       string
	Country     string
	PhoneNumber string
	Password    string
	RoleID      int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RiderProfile struct {
	ID          string
	UserID      string
	PaymentInfo string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DriverProfile struct {
	ID            string
	UserID        string
	LicenseNumber string
	VehicleNumber string
	VehicleModel  string
	VehicleMake   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type UserDetail struct {
	ID          string
	Name        string
	PhoneNumber string
	Role        string
}

type DriverDetail struct {
	UserDetail
	LicenseNumber string
	VehicleNumber string
	VehicleMake   string
	VehicleModel  string
}

type DriverVehicleMetadata struct {
	DriverID    string
	VehicleType string
}
