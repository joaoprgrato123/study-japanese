package role

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository interface {
	GetAll(ctx context.Context) ([]Role, error)
	GetByID(ctx context.Context, id int64) (Role, error)
	GetByName(ctx context.Context, name string) (Role, error)
	Create(ctx context.Context, name string) (Role, error)
	AssignToUser(ctx context.Context, userID, roleID int64) error
	RemoveFromUser(ctx context.Context, userID, roleID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]Role, error)
	UserHasRole(ctx context.Context, userID int64, roleName string) (bool, error)
}

type roleRepository struct {
	pool *pgxpool.Pool
}

func NewRoleRepository(pool *pgxpool.Pool) RoleRepository {
	return &roleRepository{pool: pool}
}

func (r *roleRepository) GetAll(ctx context.Context) ([]Role, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, name, created_at FROM roles ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("GetAll: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *roleRepository) GetByID(ctx context.Context, id int64) (Role, error) {
	var role Role
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, created_at FROM roles WHERE id = $1`,
		id,
	).Scan(&role.ID, &role.Name, &role.CreatedAt)
	if err != nil {
		return Role{}, fmt.Errorf("GetByID: %w", err)
	}
	return role, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (Role, error) {
	var role Role
	err := r.pool.QueryRow(ctx,
		`SELECT id, name, created_at FROM roles WHERE name = $1`,
		name,
	).Scan(&role.ID, &role.Name, &role.CreatedAt)
	if err != nil {
		return Role{}, fmt.Errorf("GetByName: %w", err)
	}
	return role, nil
}

func (r *roleRepository) Create(ctx context.Context, name string) (Role, error) {
	var role Role
	err := r.pool.QueryRow(ctx,
		`INSERT INTO roles (name) VALUES ($1) RETURNING id, name, created_at`,
		name,
	).Scan(&role.ID, &role.Name, &role.CreatedAt)
	if err != nil {
		return Role{}, fmt.Errorf("Create: %w", err)
	}
	return role, nil
}

func (r *roleRepository) AssignToUser(ctx context.Context, userID, roleID int64) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2)
		 ON CONFLICT DO NOTHING`,
		userID, roleID,
	)
	if err != nil {
		return fmt.Errorf("AssignToUser: %w", err)
	}
	return nil
}

func (r *roleRepository) RemoveFromUser(ctx context.Context, userID, roleID int64) error {
	_, err := r.pool.Exec(ctx,
		`DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`,
		userID, roleID,
	)
	if err != nil {
		return fmt.Errorf("RemoveFromUser: %w", err)
	}
	return nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID int64) ([]Role, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT r.id, r.name, r.created_at
		 FROM roles r
		 INNER JOIN user_roles ur ON ur.role_id = r.id
		 WHERE ur.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("GetUserRoles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name, &role.CreatedAt); err != nil {
			return nil, fmt.Errorf("scanning role: %w", err)
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (r *roleRepository) UserHasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM user_roles ur
			INNER JOIN roles r ON r.id = ur.role_id
			WHERE ur.user_id = $1 AND r.name = $2
		)`,
		userID, roleName,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("UserHasRole: %w", err)
	}
	return exists, nil
}
