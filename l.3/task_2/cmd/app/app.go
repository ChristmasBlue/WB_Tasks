package app

import (
	"fmt"
	"log"

	"task-2/internal/cache/redis"
	"task-2/internal/config"
	"task-2/internal/handler"
	"task-2/internal/repository"
	"task-2/internal/service"

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
	cache := redis.New()
	service := service.New(repository, cache)
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
	engine.POST("/shorten", handler.CreateShortUrl)

	// GET requests
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	engine.GET("/", handler.GetMainPage)
	engine.GET("/s/:short_url", handler.RedirectByShortUrl)
	engine.GET("analytics/:short_url", handler.GetAnalytics)
	engine.GET("analytics/user_agent", handler.AggregateByUserAgent)
	engine.GET("analytics/date", handler.AggregateByDate)
	engine.GET("analytics/month", handler.AggregateByMonth)
}
