package domain

type DriverNotifier interface {
	NotifyClient(clientID string, payload any) error
}
