package router

import (
	"database/sql"
	"net/http"
	"newww/internal/apps/category"
)

func registerCategoryRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := category.NewRepository(db)
	service := category.NewService(repo)
	handler := category.NewHandler(service)

	mux.HandleFunc("/api/categories", handler.List)
}
