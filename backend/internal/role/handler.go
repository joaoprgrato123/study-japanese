package role

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joaoprgrato123/study-japanese/internal/httputil"
)

type RoleHandler struct {
	service RoleService
}

func NewRoleHandler(service RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

func (h *RoleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	roles, err := h.service.GetAll(r.Context())
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toRoleResponses(roles))
}

func (h *RoleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var body CreateRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	role, err := h.service.Create(r.Context(), body.Name)
	if err != nil {
		switch {
		case errors.Is(err, ErrRoleNameEmpty):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "role name cannot be empty"})
		case errors.Is(err, ErrRoleNameTaken):
			httputil.WriteJSON(w, http.StatusConflict, map[string]string{"error": "role name already exists"})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, toRoleResponse(role))
}

func (h *RoleHandler) AssignToUser(w http.ResponseWriter, r *http.Request) {
	var body AssignRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.AssignToUser(r.Context(), body.UserID, body.RoleID); err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			httputil.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "role not found"})
			return
		}
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "role assigned"})
}

func (h *RoleHandler) RemoveFromUser(w http.ResponseWriter, r *http.Request) {
	var body AssignRoleRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	if err := h.service.RemoveFromUser(r.Context(), body.UserID, body.RoleID); err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			httputil.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "role not found"})
			return
		}
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "role removed"})
}

func (h *RoleHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(mux.Vars(r)["user_id"], 10, 64)
	if err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid user id"})
		return
	}

	roles, err := h.service.GetUserRoles(r.Context(), userID)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, toRoleResponses(roles))
}
