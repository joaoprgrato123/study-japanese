package user

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetByID(ctx context.Context, id int64) (User, error)
	GetByEmail(ctx context.Context, email string) (User, error)
	Create(ctx context.Context, email, hashedPassword string) (User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) GetByID(ctx context.Context, id int64) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return User{}, fmt.Errorf("GetByID: %w", err)
	}
	return u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return User{}, fmt.Errorf("GetByEmail: %w", err)
	}
	return u, nil
}

func (r *userRepository) Create(ctx context.Context, email, hashedPassword string) (User, error) {
	var u User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password)
		 VALUES ($1, $2)
		 RETURNING id, email, created_at`,
		email, hashedPassword,
	).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		return User{}, fmt.Errorf("Create: %w", err)
	}
	return u, nil
}
