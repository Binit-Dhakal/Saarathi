package domain

const (
	RoleAdmin = iota + 1
	RoleRider
	RoleDriver
)

type Permission string

type Roles struct {
	ID          string
	Name        string
	Permissions []Permission
}
