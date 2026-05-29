package main

import (
	"log"

	_ "task-3/docs"

	"task-3/cmd/app"
)

// @title Comment Tree API
// @version 1.0
// @description A Comment Tree API where can work with comment tree.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	if err := app.Run(); err != nil {
		log.Fatal("could not start server: ", err)
	}
}
