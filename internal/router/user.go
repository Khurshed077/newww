package router

import (
	"database/sql"
	"net/http"

	"newww/internal/apps/user"
	"newww/internal/middleware"
)

func registerUserRoutes(mux *http.ServeMux, db *sql.DB) {
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewUserHandler(userService)

	mux.Handle(
		"/users",
		middleware.Middleware(http.HandlerFunc(userHandler.Users)),
	)
}
