package handler

import (
	"errors"
	"net/http"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	_ "task-2/internal/dto"
	"task-2/internal/model"
	"task-2/internal/repository"
)

// RedirectByShortUrl godoc
// @Summary Redirect to original URL by short URL
// @Description Redirects to the original URL corresponding to the given short URL
// @Tags URL
// @Produce plain
// @Param short_url path string true "Short URL"
// @Success 301 "Redirect to original URL"
// @Failure 400 {object} map[string]string "Invalid short URL or not found"
// @Router /s/{short_url} [get]
func (h *Handler) RedirectByShortUrl(c *ginext.Context) {
	short_url := c.Param("short_url")
	var redirectInfo model.RedirectInfo
	redirectInfo.ShortUrl = short_url
	redirectInfo.UserAgent = c.Request.UserAgent()

	url, err := h.service.GetUrlByShort(short_url, redirectInfo)
	if err != nil {
		zlog.Logger.Error().Msg("could not get short url: " + err.Error())
		if errors.Is(err, repository.ErrAliasNotFound) {
			c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("successfully handled GET request: " + url.Url)
	c.Redirect(http.StatusMovedPermanently, url.Url)
}

// GetAnalytics godoc
// @Summary Get analytics data for a short URL
// @Description Returns analytics data for the given short URL
// @Tags Analytics
// @Produce json
// @Param short_url path string true "Short URL"
// @Success 200 {object} dto.RedirectInfo
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /analytics/{short_url} [get]
func (h *Handler) GetAnalytics(c *ginext.Context) {
	short_url := c.Param("short_url")
	analytics, err := h.service.GetAnalytics(short_url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	zlog.Logger.Info().Msg("succesfully handled GET request for geting analytics data")
	c.JSON(http.StatusOK, analytics)
}

// GetMainPage godoc
// @Summary Get main page
// @Description Returns the main HTML page of the URL shortener
// @Tags UI
// @Produce html
// @Success 200 "HTML page"
// @Router / [get]
func (h *Handler) GetMainPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
