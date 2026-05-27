package main

import (
	"log"

	_ "task-1/docs"

	"task-1/cmd/app"
)

// @title Delayed Notifier API
// @version 1.0
// @description A Delayed notifier which plans noitfication and send on needed time.
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
