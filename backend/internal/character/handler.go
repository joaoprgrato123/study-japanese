package character

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/joaoprgrato123/study-japanese/internal/httputil"
)

type CharacterHandler struct {
	service CharacterService
}

func NewCharacterHandler(service CharacterService) *CharacterHandler {
	return &CharacterHandler{service: service}
}

// GenerateChallenge is a protected route — any logged in user can start a challenge
func (h *CharacterHandler) GenerateChallenge(w http.ResponseWriter, r *http.Request) {
	var req ChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	challenge, err := h.service.GenerateChallenge(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCategory):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid category"})
		case errors.Is(err, ErrInvalidLevel):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid level"})
		case errors.Is(err, ErrInvalidAmount):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "amount must be greater than 0"})
		case errors.Is(err, ErrNotEnoughChars):
			httputil.WriteJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": "not enough characters for this filter"})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusOK, challenge)
}

// AddCharacter is a protected route — only admins should call this (roles come later)
func (h *CharacterHandler) AddCharacter(w http.ResponseWriter, r *http.Request) {
	var c Character
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}

	created, err := h.service.AddCharacter(r.Context(), c)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidCategory):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid category"})
		case errors.Is(err, ErrInvalidLevel):
			httputil.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid level"})
		default:
			httputil.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "something went wrong"})
		}
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, created)
}
