package router

import (
	"net/http"
	_ "newww/docs" // <- путь к swag doc

	httpSwagger "github.com/swaggo/http-swagger"
)

func RegisterSwagger(mux *http.ServeMux) {
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
}
