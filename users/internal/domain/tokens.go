package domain

import "time"

type Token struct {
	UserID       string    `json:"userId"`
	RefreshToken string    `json:"refreshToken"`
	RoleID       int       `json:"roleId"`
	ExpiresAt    time.Time `json:"expiresAt"`
}
