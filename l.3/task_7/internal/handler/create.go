package handler

import (
	"net/http"
	"task-7/internal/auth"
	"task-7/internal/config"
	"task-7/internal/dto"
	_ "task-7/internal/model"
	"task-7/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// @Summary Create a new user
// @Description Create a new user with name and role, returns user with JWT token
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.CreateUser true "User creation data"
// @Success 200 {object} model.User
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /users [post]
func (h *Handler) CreateUser(c *ginext.Context) {
	var createUser dto.CreateUser

	if err := c.BindJSON(&createUser); err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(createUser); err != nil {
		errMsg := validator.CreateValidationErrorResponse(err)
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	user, err := h.service.CreateUser(h.ctx, createUser)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	secret := []byte(config.Cfg.HttpServer.Secret)
	token, err := auth.CreateJWT(secret, user.Role)
	if err != nil {
		zlog.Logger.Error().Msg("could not create jwt token: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not create jwt token"})
		return
	}
	user.Token = token

	zlog.Logger.Info().Msg("successfully handled request and created new user")
	c.JSON(http.StatusOK, user)
}

// @Summary Create a new item
// @Description Create a new item, requires admin authentication
// @Tags items
// @Accept json
// @Produce json
// @Param item body dto.CreateItem true "Item creation data"
// @Success 200 {object} model.Item
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /items [post]
// @Security BearerAuth
func (h *Handler) CreateItem(c *ginext.Context) {
	var createItem dto.CreateItem

	if err := c.BindJSON(&createItem); err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(createItem); err != nil {
		errMsg := validator.CreateValidationErrorResponse(err)
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	item, err := h.service.CreateItem(h.ctx, createItem)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled request and created new item")
	c.JSON(http.StatusOK, item)
}
