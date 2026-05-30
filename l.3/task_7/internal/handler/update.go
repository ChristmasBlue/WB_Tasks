package handler

import (
	"net/http"
	"strconv"
	"task-7/internal/dto"
	"task-7/internal/validator"

	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// @Summary Update an item
// @Description Update an item by ID, requires admin or manager authentication
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param item body dto.UpdateItem true "Update data"
// @Success 200 {object} gin.H
// @Failure 400 {object} gin.H
// @Failure 500 {object} gin.H
// @Router /items/{id} [put]
// @Security BearerAuth
func (h *Handler) UpdateItem(c *ginext.Context) {
	id := c.Param("id")
	itemID, err := strconv.Atoi(id)
	if err != nil {
		zlog.Logger.Error().Msg("invalid id: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item is or was not provided"})
		return
	}

	var updateItem dto.UpdateItem
	updateItem.ID = itemID
	if err := c.BindJSON(&updateItem); err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validator.Validate.Struct(updateItem); err != nil {
		errMsg := validator.CreateValidationErrorResponse(err)
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
		return
	}

	err = h.service.UpdateItem(h.ctx, &updateItem)
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled request and updated item")
	c.JSON(http.StatusOK, gin.H{"status": "successfully updated item"})
}
