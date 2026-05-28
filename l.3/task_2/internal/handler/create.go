package handler

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"task-2/internal/model"
)

// CreateShortUrl godoc
// @Summary Create a shortened URL
// @Description Create a shortened URL from the provided original URL
// @Tags URL
// @Accept json
// @Produce json
// @Param url body model.Url true "URL to shorten"
// @Success 200 {object} model.Url
// @Failure 400 {object} ginext.H "Invalid request body"
// @Failure 500 {object} ginext.H "Internal server error"
// @Router /shorten [post]
func (h *Handler) CreateShortUrl(c *ginext.Context) {
	var url model.Url
	if err := c.BindJSON(&url); err != nil {
		zlog.Logger.Error().Msg("could not bind json to object: " + err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid request body",
		})
		return
	}

	urlInfo, err := h.service.CreateShortUrl(url)
	if err != nil {
		zlog.Logger.Error().Msg("could not create short_url: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
	}

	zlog.Logger.Info().Msg("successfully handled GET request and created short url for url")
	c.JSON(http.StatusOK, urlInfo)
}
