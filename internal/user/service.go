package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailTaken         = errors.New("email already in use")
	ErrInvalidInput       = errors.New("invalid input")
)

type UserService interface {
	Register(ctx context.Context, email, password string) (User, error)
	Login(ctx context.Context, email, password string) (User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(ctx context.Context, email, password string) (User, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	if email == "" || password == "" {
		return User{}, ErrInvalidInput
	}
	if len(password) < 8 {
		return User{}, fmt.Errorf("%w: password must be at least 8 characters", ErrInvalidInput)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("hashing password: %w", err)
	}

	u, err := s.repo.Create(ctx, email, string(hashed))
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return User{}, ErrEmailTaken
		}
		return User{}, fmt.Errorf("creating user: %w", err)
	}

	return u, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (User, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	u, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		// don't reveal whether the email exists or not
		// always return the same error regardless
		return User{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return User{}, ErrInvalidCredentials
	}

	return u, nil
}
