package domain

type LocationRepo interface {
	SaveActiveGeoLocation(*DriverLocation) error
	RemoveActiveGeoLocation(string) error
}
