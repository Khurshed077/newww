package app

import (
	"log"
	"net/http"
	"newww/internal/middleware"
	"newww/internal/router"
	"newww/internal/stroge"
)

func Run() {
	db := stroge.NewSQLite("data.db")
	defer db.Close()
	mux := router.New(db)

	handler := middleware.CorsMiddleware(mux)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
