package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"task-5/internal/config"
	"task-5/internal/handler"
	"task-5/internal/rabbitmq"
	"task-5/internal/repository"
	"task-5/internal/sender"
	"task-5/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

const (
	MaxOpenConns = 10
	MaxIdleConns = 5
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

	opts := &dbpg.Options{MaxOpenConns: MaxOpenConns, MaxIdleConns: MaxIdleConns}
	db, err := dbpg.New(dbString, []string{}, opts)
	if err != nil {
		log.Fatal("could not init db: " + err.Error())
	}

	repository := repository.New(db)
	queue := rabbitmq.New()
	defer queue.Close() // Закрываем соединение при завершении

	sender := sender.New()
	svc := service.New(repository, queue, sender)
	handler := handler.New(svc)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("received shutting signal %v. Shutting down", sig)
		cancel()
	}()

	svc.StartWorker(ctx)

	router := ginext.New("release")
	registerRoutes(router, handler)

	zlog.Logger.Info().Msg("successfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	engine.LoadHTMLFiles("/app/static/index.html", "/app/static/admin.html", "/app/static/user.html")
	engine.Static("/static", "/app/static")

	engine.POST("/create_event", handler.CreateEvent)
	engine.POST("/events/:id/book", handler.CreateBooking)
	engine.POST("/events/:id/confirm", handler.ConfirmPayment)

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/admin", handler.GetAdminPage)
	engine.GET("/user", handler.GetUserPage)

	engine.GET("/events/:id", handler.GetEvent)
	engine.GET("/events", handler.GetAllEvents)
}
