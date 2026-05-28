package handler

import (
	"net/http"

	_ "task-2/internal/dto"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

// AggregateByUserAgent godoc
// @Summary Get aggregated analytics by user agent
// @Description Returns aggregated analytics data grouped by user agent
// @Tags Analytics
// @Produce json
// @Success 200 {array} dto.UserAgentDTO
// @Failure 500 {object} ginext.H "Internal server error"
// @Router /analytics/user_agent [get]
func (h *Handler) AggregateByUserAgent(c *ginext.Context) {
	analytics, err := h.service.AggregateByUserAgent()
	if err != nil {
		zlog.Logger.Error().Msg("could not get aggregated data by user agent from db: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get aggregated data by user agent from db: " + err.Error(),
		})
		return
	}

	zlog.Logger.Info().Msg("succesfully handled GET request for getting aggreagated by user_agent data")
	c.JSON(http.StatusOK, analytics)
}

// AggregateByDate godoc
// @Summary Get aggregated analytics by date
// @Description Returns aggregated analytics data grouped by date
// @Tags Analytics
// @Produce json
// @Success 200 {array} dto.DateDTO
// @Failure 500 {object} ginext.H "Internal server error"
// @Router /analytics/date [get]
func (h *Handler) AggregateByDate(c *ginext.Context) {
	analytics, err := h.service.AggregateByDate()
	if err != nil {
		zlog.Logger.Error().Msg("could not get aggregated data by date from db: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get aggregated data by date from db: " + err.Error(),
		})
		return
	}

	zlog.Logger.Info().Msg("succesfully handled GET request for getting aggreagated by date data")
	c.JSON(http.StatusOK, analytics)
}

// AggregateByMonth godoc
// @Summary Get aggregated analytics by month
// @Description Returns aggregated analytics data grouped by month
// @Tags Analytics
// @Produce json
// @Success 200 {array} dto.MonthDTO
// @Failure 500 {object} ginext.H "Internal server error"
// @Router /analytics/month [get]
func (h *Handler) AggregateByMonth(c *ginext.Context) {
	analytics, err := h.service.AggregateByMonth()
	if err != nil {
		zlog.Logger.Error().Msg("could not get aggregated data by month from db: " + err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get aggregated data by month from db: " + err.Error(),
		})
		return
	}

	zlog.Logger.Info().Msg("succesfully handled GET request for getting aggreagated by month data")
	c.JSON(http.StatusOK, analytics)
}
