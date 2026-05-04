package main

import (
	"task-18/internal/config"
	"task-18/internal/di"
	"task-18/internal/logger"
	"task-18/internal/repository"
	"task-18/internal/web"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			config.MustLoad,
			logger.ProvideLogger,
			repository.NewInMemoryRepo,
			func(repo *repository.InMemoryRepo) repository.Storage {
				return repo
			},
			web.NewCalendarHandler,
		),

		fx.Invoke(
			di.StartHttpServer,
		),
	)
	app.Run()
}
