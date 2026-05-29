package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "task-4/docs"
	"task-4/internal/config"
	"task-4/internal/handler"
	"task-4/internal/kafka"
	repository "task-4/internal/repository/db"
	"task-4/internal/repository/storage/minio"
	"task-4/internal/service"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
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

	repository := repository.NewPostgres(db)
	queue := kafka.New()
	fileStorage := minio.New()
	service := service.New(repository, fileStorage, queue)
	handler := handler.New(service)

	router := ginext.New("release")
	registerRoutes(router, handler)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting signal %v. Shuting down", sig)
		cancel()
	}()

	service.StartWorkers(ctx)

	zlog.Logger.Info().Msg("succesfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	// Register static files
	engine.LoadHTMLFiles("static/index.html")
	engine.Static("/static", "static")

	// POST requests
	engine.POST("/upload", handler.CreateImage)

	// GET requests
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/image/:id", handler.GetImageByID)
	engine.GET("/image/info/:id", handler.GetImageInfo)

	// DELETE request
	engine.DELETE("/image/:id", handler.DeleteImageByID)
}
