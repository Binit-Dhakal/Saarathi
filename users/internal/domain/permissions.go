package domain

type Permission string

var (
	PermissionAcceptRide Permission = "accept_ride"
	PermissionCheckFare  Permission = "check_fare"
	PermissionFullAccess Permission = "full_access"
)
