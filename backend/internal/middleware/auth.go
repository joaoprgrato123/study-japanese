package middleware

import (
	"context"
	"net/http"

	"github.com/joaoprgrato123/study-japanese/internal/role"
	"github.com/joaoprgrato123/study-japanese/internal/session"
)

// unexported type prevents collision with context keys from other packages
type contextKey string

const SessionContextKey contextKey = "session"

func RequireAuth(sessionService session.SessionService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			sess, err := sessionService.GetByToken(r.Context(), cookie.Value)
			if err != nil {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}

			// store session in context so downstream handlers can access it
			ctx := context.WithValue(r.Context(), SessionContextKey, sess)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roleService role.RoleService, roleName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sess := r.Context().Value(SessionContextKey).(session.Session)

			hasRole, err := roleService.UserHasRole(r.Context(), sess.UserID, roleName)
			if err != nil || !hasRole {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
