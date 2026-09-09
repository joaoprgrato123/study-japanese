package role

import "time"

type Role struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type AssignRoleRequest struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}
