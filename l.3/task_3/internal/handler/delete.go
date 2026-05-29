package handler

import (
	"net/http"
	"strconv"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// @Summary Удалить комментарий по ID
// @Description Удаляет комментарий по его уникальному идентификатору
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "ID комментария для удаления"
// @Success 200 {object} map[string]string "status":"successfully deleted comment"
// @Failure 400 {object} map[string]string "error":"invalid id was provided" or "could not delete comment from db"
// @Router /comments/{id} [delete]
func (h *Handler) DeleteCommentById(c *ginext.Context) {
	id := c.Param("id")
	commentId, err := strconv.Atoi(id)
	if err != nil {
		zlog.Logger.Error().Msgf("invalid id was provided: %s", err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid id was provided",
		})
		return
	}

	if err := h.service.DeleteCommentById(commentId); err != nil {
		zlog.Logger.Error().Msgf("could not delete comment from db: %s", err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "could not delete comment from db: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, ginext.H{
		"status": "successfully deleted comment",
	})
}
