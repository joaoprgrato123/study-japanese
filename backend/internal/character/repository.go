package character

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CharacterRepository interface {
	GetRandom(ctx context.Context, category, level string, amount int) ([]Character, error)
	GetDistractors(ctx context.Context, excludeIDs []int64, amount int) ([]Character, error)
	Create(ctx context.Context, c Character) (Character, error)
}

type characterRepository struct {
	pool *pgxpool.Pool
}

func NewCharacterRepository(pool *pgxpool.Pool) CharacterRepository {
	return &characterRepository{pool: pool}
}

// GetRandom fetches N random characters filtered by category and level.
// If category is "random" it skips the category filter.
func (r *characterRepository) GetRandom(ctx context.Context, category, level string, amount int) ([]Character, error) {
	var rows pgx.Rows
	var err error

	if category == "random" {
		rows, err = r.pool.Query(ctx,
			`SELECT id, character, reading, meaning, category, level, created_at
			 FROM characters
			 WHERE level = $1
			 ORDER BY RANDOM()
			 LIMIT $2`,
			level, amount,
		)
	} else {
		rows, err = r.pool.Query(ctx,
			`SELECT id, character, reading, meaning, category, level, created_at
			 FROM characters
			 WHERE category = $1 AND level = $2
			 ORDER BY RANDOM()
			 LIMIT $3`,
			category, level, amount,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("GetRandom: %w", err)
	}
	defer rows.Close()

	return scanRows(rows)
}

// GetDistractors fetches random characters excluding the ones already used as questions.
// These become the wrong options in the multiple choice.
func (r *characterRepository) GetDistractors(ctx context.Context, excludeIDs []int64, amount int) ([]Character, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, character, reading, meaning, category, level, created_at
		 FROM characters
		 WHERE id != ALL($1)
		 ORDER BY RANDOM()
		 LIMIT $2`,
		excludeIDs, amount,
	)
	if err != nil {
		return nil, fmt.Errorf("GetDistractors: %w", err)
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *characterRepository) Create(ctx context.Context, c Character) (Character, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO characters (character, reading, meaning, category, level)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at`,
		c.Character, c.Reading, c.Meaning, c.Category, c.Level,
	).Scan(&c.ID, &c.CreatedAt)
	if err != nil {
		return Character{}, fmt.Errorf("Create: %w", err)
	}
	return c, nil
}

// scanRows is a private helper to avoid repeating scan logic
func scanRows(rows pgx.Rows) ([]Character, error) {
	var characters []Character
	for rows.Next() {
		var c Character
		if err := rows.Scan(&c.ID, &c.Character, &c.Reading, &c.Meaning, &c.Category, &c.Level, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning character: %w", err)
		}
		characters = append(characters, c)
	}
	return characters, rows.Err()
}
