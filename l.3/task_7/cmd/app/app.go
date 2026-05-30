package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"task-7/internal/config"
	"task-7/internal/handler"
	"task-7/internal/middleware"
	"task-7/internal/repository"
	"task-7/internal/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

func Run() error {
	zlog.Init()

	gin.SetMode(gin.ReleaseMode)

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
	service := service.New(repository)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		zlog.Logger.Info().Msgf("recieved shutting signal %v. Shuting down", sig)
		cancel()
	}()

	router := gin.New()
	handler := handler.New(ctx, service)
	registerRoutes(router, handler)

	zlog.Logger.Info().Msg("succesfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *gin.Engine, handler *handler.Handler) {
	engine.LoadHTMLFiles("/app/static/login.html", "/app/static/main.html")
	engine.Static("/static", "/app/static")

	group := engine.Group("/", middleware.AuthMiddleware([]byte(config.Cfg.HttpServer.Secret))) // Register static files

	// POST requests
	group.POST("/items", handler.CreateItem)
	engine.POST("/users", handler.CreateUser)

	// GET requests
	engine.GET("/login", handler.GetLoginPage)
	engine.GET("/main", handler.GetMainPage)
	group.GET("/items", handler.GetAllItems)
	group.GET("/users/history", handler.GetUsersWithChanges)

	// PUT requests
	group.PUT("/items/:id", handler.UpdateItem)

	// DELETE requests
	group.DELETE("/items/:id", handler.DeleteItem)

	// swagger
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
