package category

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// @Summary List categories
// @Description Get all categories
// @Tags Categories
// @Produce json
// @Success 200 {array} model.Category
// @Router /api/categories [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
