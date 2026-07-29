package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joaoprgrato123/study-japanese/internal/session"
)

type UserHandler struct {
	userService    UserService
	sessionService session.SessionService
}

func NewUserHandler(userService UserService, sessionService session.SessionService) *UserHandler {
	return &UserHandler{userService: userService, sessionService: sessionService}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	u, err := h.userService.Register(r.Context(), body.Email, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			writeJSON(w, http.StatusConflict, map[string]string{"error": "email already in use"})
		case errors.Is(err, ErrInvalidInput):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		}
		return
	}

	writeJSON(w, http.StatusCreated, UserResponse{ID: u.ID, Email: u.Email})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	u, err := h.userService.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	sess, err := h.sessionService.Create(r.Context(), u.ID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sess.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	writeJSON(w, http.StatusOK, UserResponse{ID: u.ID, Email: u.Email})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "not logged in"})
		return
	}

	if err := h.sessionService.Delete(r.Context(), cookie.Value); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
