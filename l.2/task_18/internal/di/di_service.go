package di

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"log"
	"net/http"
	"task-18/internal/config"
	"task-18/internal/web"
)

func StartHttpServer(lc fx.Lifecycle, calendarHandler *web.CalendarHandler, config *config.Config) {
	router := chi.NewRouter()

	web.RegisterRoutes(router, calendarHandler)
	address := fmt.Sprintf(":%d", config.HttpPort)
	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Server started")
			go func() {
				if err := server.ListenAndServe(); err != nil {
					log.Printf("ListenAndServe error: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Printf("Shutting down server...")
			return server.Close()
		},
	})

}
