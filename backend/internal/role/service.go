package role

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrRoleNotFound  = errors.New("role not found")
	ErrRoleNameEmpty = errors.New("role name cannot be empty")
	ErrRoleNameTaken = errors.New("role name already exists")
)

type RoleService interface {
	GetAll(ctx context.Context) ([]Role, error)
	Create(ctx context.Context, name string) (Role, error)
	AssignToUser(ctx context.Context, userID, roleID int64) error
	RemoveFromUser(ctx context.Context, userID, roleID int64) error
	GetUserRoles(ctx context.Context, userID int64) ([]Role, error)
	UserHasRole(ctx context.Context, userID int64, roleName string) (bool, error)
}

type roleService struct {
	repo RoleRepository
}

func NewRoleService(repo RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) GetAll(ctx context.Context) ([]Role, error) {
	return s.repo.GetAll(ctx)
}

func (s *roleService) Create(ctx context.Context, name string) (Role, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return Role{}, ErrRoleNameEmpty
	}

	role, err := s.repo.Create(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), "unique") {
			return Role{}, ErrRoleNameTaken
		}
		return Role{}, fmt.Errorf("creating role: %w", err)
	}

	return role, nil
}

func (s *roleService) AssignToUser(ctx context.Context, userID, roleID int64) error {
	// verify the role exists before assigning
	_, err := s.repo.GetByID(ctx, roleID)
	if err != nil {
		return ErrRoleNotFound
	}

	return s.repo.AssignToUser(ctx, userID, roleID)
}

func (s *roleService) RemoveFromUser(ctx context.Context, userID, roleID int64) error {
	_, err := s.repo.GetByID(ctx, roleID)
	if err != nil {
		return ErrRoleNotFound
	}

	return s.repo.RemoveFromUser(ctx, userID, roleID)
}

func (s *roleService) GetUserRoles(ctx context.Context, userID int64) ([]Role, error) {
	return s.repo.GetUserRoles(ctx, userID)
}

func (s *roleService) UserHasRole(ctx context.Context, userID int64, roleName string) (bool, error) {
	return s.repo.UserHasRole(ctx, userID, roleName)
}
