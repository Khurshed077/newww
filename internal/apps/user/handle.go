package user

import (
	"encoding/json"
	"net/http"
	"newww/internal/middleware/jwt"
)

type UserHandler struct {
	service *UserService
}

func NewUserHandler(s *UserService) *UserHandler {
	return &UserHandler{service: s}
}
func (h *UserHandler) Users(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(jwt.ClaimsKey).(*jwt.Claims)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetByID(claims.UserID)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
