package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type SessionService interface {
	Create(ctx context.Context, userID int64) (Session, error)
	GetByToken(ctx context.Context, token string) (Session, error)
	Delete(ctx context.Context, token string) error
}

type sessionService struct {
	repo SessionRepository
}

func NewSessionService(repo SessionRepository) SessionService {
	return &sessionService{repo: repo}
}

func (s *sessionService) Create(ctx context.Context, userID int64) (Session, error) {
	token, err := generateToken()
	if err != nil {
		return Session{}, fmt.Errorf("generating token: %w", err)
	}

	expiresAt := time.Now().Add(7 * 24 * time.Hour)

	return s.repo.Create(ctx, userID, token, expiresAt)
}

func (s *sessionService) GetByToken(ctx context.Context, token string) (Session, error) {
	return s.repo.GetByToken(ctx, token)
}

func (s *sessionService) Delete(ctx context.Context, token string) error {
	return s.repo.Delete(ctx, token)
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
