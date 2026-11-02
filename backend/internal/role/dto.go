package role

import "time"

// --- requests ---

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type AssignRoleRequest struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}

// --- responses ---

type RoleResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// --- mappers ---

func toRoleResponse(r Role) RoleResponse {
	return RoleResponse{
		ID:        r.ID,
		Name:      r.Name,
		CreatedAt: r.CreatedAt,
	}
}

func toRoleResponses(roles []Role) []RoleResponse {
	out := make([]RoleResponse, 0, len(roles))
	for _, r := range roles {
		out = append(out, toRoleResponse(r))
	}
	return out
}
