package domain

type TripStatus string

var (
	TripStatusPending   TripStatus = "pending"
	TripStatusApproved  TripStatus = "approved"
	TripStatusCancelled TripStatus = "cancelled"
)

func (t TripStatus) String() string {
	switch t {
	case TripStatusPending, TripStatusApproved, TripStatusCancelled:
		return string(t)
	default:
		return ""
	}
}
