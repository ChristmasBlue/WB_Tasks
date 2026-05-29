package app

import (
	"fmt"
	"log"

	_ "task-3/docs"

	"task-3/internal/config"
	"task-3/internal/handler"
	"task-3/internal/repository"
	"task-3/internal/service"

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

	repository := repository.New(db)
	service := service.New(repository)
	handler := handler.New(service)

	router := ginext.New("release")
	registerRoutes(router, handler)

	zlog.Logger.Info().Msg("succesfully started server on " + config.Cfg.HttpServer.Address)
	return router.Run(config.Cfg.HttpServer.Address)
}

func registerRoutes(engine *ginext.Engine, handler *handler.Handler) {
	// Register static files
	engine.LoadHTMLFiles("/app/static/index.html")
	engine.Static("/static", "/app/static")

	// POST requests
	engine.POST("/comments", handler.CreateComment)
	engine.POST("/comments/search", handler.GetCommentsByTextSearch)

	// GET requests
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/comments", handler.GetComments)
	engine.GET("/comments/all", handler.GetAllComments)

	// DELETE request
	engine.DELETE("/comments/:id", handler.DeleteCommentById)

}
