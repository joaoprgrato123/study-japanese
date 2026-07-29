package session

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository interface {
	Create(ctx context.Context, userID int64, token string, expiresAt time.Time) (Session, error)
	GetByToken(ctx context.Context, token string) (Session, error)
	Delete(ctx context.Context, token string) error
}

type sessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionRepository(pool *pgxpool.Pool) SessionRepository {
	return &sessionRepository{pool: pool}
}

func (r *sessionRepository) Create(ctx context.Context, userID int64, token string, expiresAt time.Time) (Session, error) {
	var s Session
	err := r.pool.QueryRow(ctx,
		`INSERT INTO sessions (user_id, token, expires_at)
		 VALUES ($1, $2, $3)
		 RETURNING id, user_id, token, expires_at`,
		userID, token, expiresAt,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt)
	if err != nil {
		return Session{}, fmt.Errorf("Create session: %w", err)
	}
	return s, nil
}

func (r *sessionRepository) GetByToken(ctx context.Context, token string) (Session, error) {
	var s Session
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, token, expires_at FROM sessions
		 WHERE token = $1 AND expires_at > NOW()`,
		token,
	).Scan(&s.ID, &s.UserID, &s.Token, &s.ExpiresAt)
	if err != nil {
		return Session{}, fmt.Errorf("GetByToken: %w", err)
	}
	return s, nil
}

func (r *sessionRepository) Delete(ctx context.Context, token string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
	return err
}
