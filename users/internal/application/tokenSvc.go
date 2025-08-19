package application

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/claims"
	"github.com/Binit-Dhakal/Saarathi/users/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenService interface {
	GenerateAccessAndRefreshTokens(userID string, roleID int) (*Token, error)
	ValidateRefreshToken(refreshToken string) (*domain.Token, error)
	RevokeRefreshToken(refreshToken string) error
}

type Token struct {
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
}

type JWTService struct {
	secretKey *rsa.PrivateKey
	tokenRepo domain.TokenRepo
}

func getPrivateKey(keyString string) (*rsa.PrivateKey, error) {
	// Decode the PEM encoded key
	block, _ := pem.Decode([]byte(keyString))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the key")
	}

	// Parse the key
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return key.(*rsa.PrivateKey), nil
}

func NewJWTService(secretKey string, tokenRepo domain.TokenRepo) *JWTService {
	jwtSecretKey, err := getPrivateKey(secretKey)
	if err != nil {
		panic(err)
	}
	return &JWTService{
		secretKey: jwtSecretKey,
		tokenRepo: tokenRepo,
	}
}

func (j *JWTService) GenerateAccessAndRefreshTokens(userID string, roleID int) (*Token, error) {
	permissions := []claims.Permission{}
	switch roleID {
	case domain.RoleAdmin:
		permissions = []claims.Permission{claims.PermissionFullAccess}
	case domain.RoleRider:
		permissions = []claims.Permission{claims.PermissionCheckFare}
	case domain.RoleDriver:
		permissions = []claims.Permission{claims.PermissionAcceptRide}
	}

	accessClaims := &claims.CustomClaims{
		UserID:      userID,
		RoleID:      roleID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			Issuer:    "saarathi",
			Subject:   userID,
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, accessClaims).SignedString(j.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	jti, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate UUID for JIT: %w", err)
	}

	refreshExpiry := time.Now().Add(time.Hour * 24 * 7)
	refreshClaims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshExpiry),
		Issuer:    "saarathi",
		Subject:   userID,
		ID:        jti.String(),
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims).SignedString(j.secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	token := &domain.Token{
		UserID:       userID,
		RefreshToken: refreshToken,
		RoleID:       roleID,
		ExpiresAt:    refreshExpiry,
	}

	err = j.tokenRepo.CreateToken(token)
	if err != nil {
		return nil, err
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    refreshExpiry,
	}, nil
}

func (j *JWTService) ValidateRefreshToken(refreshToken string) (*domain.Token, error) {
	token, err := j.tokenRepo.FindByRefreshToken(refreshToken)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return token, nil
}

func (j *JWTService) RevokeRefreshToken(refreshToken string) error {
	err := j.tokenRepo.RevokeRefreshToken(refreshToken)
	if err != nil {
		return err
	}

	return nil
}
