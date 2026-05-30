package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"

	_ "task-7/internal/model"
)

// @Summary Get all items
// @Description Get all items, requires authentication
// @Tags items
// @Produce json
// @Success 200 {array} model.Item
// @Failure 500 {object} gin.H
// @Router /items [get]
// @Security BearerAuth
func (h *Handler) GetAllItems(c *ginext.Context) {
	items, err := h.service.GetAllItems(h.ctx)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled request and returned all items")
	c.JSON(http.StatusOK, items)
}

// @Summary Get users with changes
// @Description Get users with their change history, requires authentication
// @Tags users
// @Produce json
// @Success 200 {array} model.UserHistory
// @Failure 500 {object} gin.H
// @Router /users/history [get]
// @Security BearerAuth
func (h *Handler) GetUsersWithChanges(c *ginext.Context) {
	users, err := h.service.GetUsersWithChanges(h.ctx)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled request and returned all items")
	c.JSON(http.StatusOK, users)
}

// GetLoginPage godoc
// @Summary      Get user page
// @Description  Get the login HTML page of the application
// @Tags         pages
// @Accept       json
// @Produce      html
// @Success      200  {string} string "HTML page content"
// @Router       / [get]
func (h *Handler) GetLoginPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

// GetMainPage godoc
// @Summary      Get main page
// @Description  Get the main HTML page of the application
// @Tags         pages
// @Accept       json
// @Produce      html
// @Success      200  {string} string "HTML page content"
// @Router       / [get]
func (h *Handler) GetMainPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "main.html", nil)
}
