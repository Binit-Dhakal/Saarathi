package domain

type LocationRepo interface {
	SaveActiveGeoLocation(*DriverLocation) error
	RemoveActiveGeoLocation(string) error
}

type WSRepo interface {
	SaveWSDetail(driverID string) error
	DeleteWSDetail(driverID string) error
}
