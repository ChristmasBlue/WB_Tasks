package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "task-1/docs"
	"task-1/internal/cache/redis"
	"task-1/internal/config"
	"task-1/internal/handler"
	"task-1/internal/rabbitmq"
	"task-1/internal/repository"
	"task-1/internal/sender"
	"task-1/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func Run() error {
	zlog.Init()

	dbString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Cfg.Postgres.Host,
		config.Cfg.Postgres.Port,
		config.Cfg.Postgres.User,
		config.Cfg.Postgres.Password,
		config.Cfg.Postgres.Name,
	)
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(dbString, []string{}, opts)
	if err != nil {
		log.Fatal("could not init db: " + err.Error())
	}

	repository := repository.New(db)
	cache := redis.New()
	queue := rabbitmq.New()
	defer queue.Close()

	sender := sender.New()
	svc := service.New(repository, cache, queue, sender)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("received shutting signal %v. Shutting down", sig)
		cancel()
	}()

	go func() {
		if err := svc.PublishReadyNotifications(ctx); err != nil {
			if err != context.Canceled {
				log.Printf("error while publishing notifications: %v", err)
			}
		}
	}()

	go func() {
		if err := svc.ConsumeMessages(ctx); err != nil {
			if err != context.Canceled {
				log.Printf("could not start consumer: %v", err)
			}
		}
	}()

	handler := handler.New(svc)

	router := ginext.New("release")
	registerRoutes(router, handler)

	zlog.Logger.Info().Msg("successfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	engine.LoadHTMLFiles("/app/static/index.html")
	engine.Static("/static", "/app/static")

	engine.POST("/notify", handler.CreateNotification)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/notify/:id", handler.GetNotificationStatus)
	engine.GET("/notify", handler.GetAllNotifications)

	engine.DELETE("notify/:id", handler.UpdateNotificationStatus)
}
