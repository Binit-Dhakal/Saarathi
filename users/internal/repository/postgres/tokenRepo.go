package postgres

import (
	"context"
	"time"

	"github.com/Binit-Dhakal/Saarathi/users/internal/domain"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TokenRepo struct {
	pool *pgxpool.Pool
}

func NewTokenRepo(pool *pgxpool.Pool) *TokenRepo {
	return &TokenRepo{
		pool: pool,
	}
}

func (t *TokenRepo) CreateToken(token *domain.Token) error {
	query := `INSERT into tokens(user_id, refresh_token,role_id, expires_at) VALUES($1,$2,$3,$4)`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{token.UserID, token.RefreshToken, token.RoleID, token.ExpiresAt}
	_, err := t.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (t *TokenRepo) FindByRefreshToken(refreshToken string) (*domain.Token, error) {
	query := `SELECT user_id, refresh_token, role_id, expires_at from tokens where refresh_token=$1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var token domain.Token
	var userUUID pgtype.UUID
	err := t.pool.QueryRow(ctx, query, refreshToken).Scan(&userUUID, &token.RefreshToken, &token.RoleID, &token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	token.UserID = userUUID.String()

	return &token, nil
}

func (t *TokenRepo) RevokeRefreshToken(refreshToken string) error {
	query := `DELETE from tokens where refresh_token=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.pool.Exec(ctx, query, refreshToken)
	if err != nil {
		return err
	}

	return nil
}
