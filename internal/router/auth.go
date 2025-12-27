package router

import (
	"database/sql"
	"html/template"
	"net/http"
	"newww/internal/apps/auth"
)

func registerAuthRoutes(mux *http.ServeMux, db *sql.DB) {
	service := auth.NewAuthService(db)
	handler := auth.NewAuthHandler(service)

	// Auth HTML
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/auth/login.html"))
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/auth/register.html"))
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/auth/logout.html"))
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/api/register", handler.Register)
	mux.HandleFunc("/api/login", handler.Login)
	mux.HandleFunc("/api/logout", handler.Logout)
	mux.HandleFunc("/api/refresh", auth.RefreshToken)

}
