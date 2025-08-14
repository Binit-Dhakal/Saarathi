package events

import "time"

type RawClickEvent struct {
	URLAlias  string
	Timestamp time.Time
	UserAgent string
	IP        string
	Referrer  string
}

type EnrichedClickEvent struct {
	URLAlias  string
	Timestamp time.Time
	Browser   string
	OS        string
	Device    string
	City      string
	Country   string
	Referrer  string
}
