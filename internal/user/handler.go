package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joaoprgrato123/study-japanese/internal/httputil"
	"github.com/joaoprgrato123/study-japanese/internal/middleware"
	"github.com/joaoprgrato123/study-japanese/internal/session"
)

type UserHandler struct {
	userService    UserService
	sessionService session.SessionService
}

func NewUserHandler(userService UserService, sessionService session.SessionService) *UserHandler {
	return &UserHandler{userService: userService, sessionService: sessionService}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var body RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	u, err := h.userService.Register(r.Context(), body.Email, body.Password)
	if err != nil {
		switch {
		case errors.Is(err, ErrEmailTaken):
			httputil.WriteJSON(w, http.StatusConflict, map[string]string{"error": "email already in use"})
		case errors.Is(err, ErrInvalidInput):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, UserResponse{ID: u.ID, Email: u.Email})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var body LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	u, err := h.userService.Login(r.Context(), body.Email, body.Password)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			httputil.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
			return
		}
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	sess, err := h.sessionService.Create(r.Context(), u.ID)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sess.Token,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	httputil.WriteJSON(w, http.StatusOK, UserResponse{ID: u.ID, Email: u.Email})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		httputil.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "not logged in"})
		return
	}

	if err := h.sessionService.Delete(r.Context(), cookie.Value); err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_id",
		Value:  "",
		MaxAge: -1,
	})

	httputil.WriteJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request) {
	sess := r.Context().Value(middleware.SessionContextKey).(session.Session)

	u, err := h.userService.GetByID(r.Context(), sess.UserID)
	if err != nil {
		httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		return
	}

	httputil.WriteJSON(w, http.StatusOK, UserResponse{ID: u.ID, Email: u.Email})
}
