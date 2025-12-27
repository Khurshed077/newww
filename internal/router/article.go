package router

import (
	"database/sql"
	"html/template"
	"net/http"
	"newww/internal/apps/article"
	"newww/internal/middleware"
)

func registerArticleRoutes(mux *http.ServeMux, db *sql.DB) {
	repo := article.NewRepository(db)
	service := article.NewService(repo)
	handler := article.NewHandler(service)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/home.html"))
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/articles/detail", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/article_detail.html"))
		tmpl.Execute(w, nil)
	})

	mux.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/dashboard.html"))
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/articles/create", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/create_article.html"))
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/articles/edit", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/edit_article.html"))
		tmpl.Execute(w, nil)
	})
	mux.HandleFunc("/articles/delete", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./templates/delete.html"))
		tmpl.Execute(w, nil)
	})
	// API маршруты
	mux.Handle(
		"/api/articles",
		middleware.CorsMiddleware(
			middleware.OptionalJWT(
				http.HandlerFunc(handler.List),
			),
		),
	)
	mux.Handle(
		"/api/articles/detail",
		middleware.CorsMiddleware(
			middleware.OptionalJWT(
				http.HandlerFunc(handler.Detail),
			),
		),
	)
	mux.Handle("/api/dashboard", middleware.JWTMiddleware(http.HandlerFunc(handler.Dashboard)))
	mux.Handle("/api/articles/create", middleware.JWTMiddleware(http.HandlerFunc(handler.Create)))
	mux.Handle("/api/articles/edit", middleware.JWTMiddleware(http.HandlerFunc(handler.Edit)))
	mux.Handle("/api/articles/delete", middleware.JWTMiddleware(http.HandlerFunc(handler.Delete)))

	// Можно оставить этот маршрут, если нужен публичный список на фронтенд
	mux.Handle(
		"/api/home",
		middleware.CorsMiddleware(
			middleware.OptionalJWT(
				http.HandlerFunc(handler.Home),
			),
		),
	)
}
