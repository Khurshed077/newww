package router

import (
	"database/sql"
	"net/http"
)

func New(db *sql.DB) *http.ServeMux {
	mux := http.NewServeMux()

	// Статика
	fs := http.FileServer(http.Dir("./uploads"))
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", fs))
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("web"))))
	// API маршруты
	RegisterSwagger(mux)
	registerArticleRoutes(mux, db) // внутри только /api/articles и т.д.
	registerAuthRoutes(mux, db)    // /api/login, /api/register, /api/logout
	registerCategoryRoutes(mux, db)
	return mux
}
