package domain

type Notifier interface {
	NotifyRider(tripID string, payload []byte)
}
