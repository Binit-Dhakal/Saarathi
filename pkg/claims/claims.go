package claims

import "github.com/golang-jwt/jwt/v5"

type Permission string

var (
	PermissionAcceptRide Permission = "accept_ride"
	PermissionCheckFare  Permission = "check_fare"
	PermissionFullAccess Permission = "full_access"
)

type CustomClaims struct {
	UserID      string
	RoleID      int
	Permissions []Permission
	jwt.RegisteredClaims
}
