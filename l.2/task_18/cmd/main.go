package main

import (
	"log"

	"task-18/cmd/api"
	"task-18/internal/config"
)

func main() {
	apiServer := api.NewServer(config.Envs.Port)
	if err := apiServer.Run(); err != nil {
		log.Fatal(err)
	}
}
