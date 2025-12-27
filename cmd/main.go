// @title Newww API
// @version 1.0
// @description API проекта Newww с JWT и ролями Admin
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	app "newww/internal"
)

func main() {

	app.Run()
}
