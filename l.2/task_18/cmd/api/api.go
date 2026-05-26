package api

import (
	"github.com/gin-gonic/gin"
	"task-18/internal/handler"
	"task-18/internal/middleware"
	"task-18/internal/repository"
	"task-18/internal/service"
)

type APIServer struct {
	addr string
}

func NewServer(addr string) *APIServer {
	return &APIServer{addr: addr}
}

func (s *APIServer) Run() error {
	router := gin.Default()
	router.Use(middleware.LoggingMiddleware()) // навесили всем хэндлерам middleware для логирования

	repository := repository.New()
	service := service.New(repository)
	handler := handler.New(service)

	router.POST("/create_event", handler.CreateEvent)
	router.POST("/update_event", handler.UpdateEvent)
	router.POST("/delete_event", handler.DeleteEvent)
	router.GET("/events_for_day", handler.GetEventsForDay)
	router.GET("/events_for_week", handler.GetEventsForWeek)
	router.GET("/events_for_month", handler.GetEventsForMonth)

	return router.Run(s.addr)
}
