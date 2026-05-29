package handler

import (
	"errors"
	"net/http"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"task-3/internal/dto"
	"task-3/internal/repository"
)

// @Summary Создать комментарий
// @Description Создает новый комментарий в системе
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body dto.CreateComment true "Данные для создания комментария"
// @Success 200 {object} dto.CreateComment "Успешно созданный комментарий"
// @Failure 400 {object} map[string]string "error":"invalid payload"
// @Failure 500 {object} map[string]string "error":"could not create comment in db"
// @Router /comments [post]
func (h *Handler) CreateComment(c *ginext.Context) {
	comment := new(dto.CreateComment)
	if err := c.BindJSON(comment); err != nil {
		zlog.Logger.Error().Msgf("could not unmarshal request body to model: %s", err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid payload: " + err.Error(),
		})
		return
	}

	comment, err := h.service.CreateComment(*comment)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidParenID) {
			zlog.Logger.Error().Msgf("could not create comment in db: %s", err.Error())
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "could not create comment in db: " + err.Error(),
			})
			return
		}

		zlog.Logger.Error().Msgf("could not create comment in db: %s", err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not create comment in db: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}
