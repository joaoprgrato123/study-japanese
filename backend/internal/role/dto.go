package role

type CreateRoleRequest struct {
	Name string `json:"name"`
}

type AssignRoleRequest struct {
	UserID int64 `json:"user_id"`
	RoleID int64 `json:"role_id"`
}
