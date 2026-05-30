package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// @Summary Delete an item
// @Description Delete an item by ID, requires admin authentication
// @Tags items
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /items/{id} [delete]
// @Security BearerAuth
func (h *Handler) DeleteItem(c *ginext.Context) {
	id := c.Param("id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item is or was not provided"})
		return
	}

	err = h.service.DeleteItem(h.ctx, itemID)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled request and deleted item")
	c.JSON(http.StatusOK, gin.H{"status": "successfully deleted item"})
}
