package domain

const (
	RoleAdmin = iota + 1
	RoleRider
	RoleDriver
)

type Roles struct {
	ID          string
	Name        string
	Permissions []Permission
}
