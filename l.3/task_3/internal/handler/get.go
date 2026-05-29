package handler

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"task-3/internal/dto"
	_ "task-3/internal/model"
)

// @Summary Получить комментарии
// @Description Получает комментарии с пагинацией и по ID родительского комментария
// @Tags comments
// @Accept json
// @Produce json
// @Param parent query int false "ID родительского комментария"
// @Param page query int false "Номер страницы для пагинации"
// @Param limit query int false "Количество комментариев на странице"
// @Success 200 {array} model.Comment "Список комментариев"
// @Failure 400 {object} map[string]string "error":"invalid payload" or "invalid id"
// @Failure 500 {object} map[string]string "error":"could not get comments"
// @Router /comments [get]
func (h *Handler) GetComments(c *ginext.Context) {
	config, err := parserQueryParameters(c.Request.URL.Query())
	if err != nil {
		zlog.Logger.Error().Msg(err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": err.Error(),
		})
		return
	}

	if config.Page == 0 && config.Limit == 0 {
		h.getCommentsById(config.ParentID, c)
		return
	}

	comments, err := h.service.GetCommentsPaginated(*config)
	if err != nil {
		zlog.Logger.Error().Msgf("could not get comments paginated: %s", err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get comments: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// @Summary Поиск комментариев по тексту
// @Description Ищет комментарии, содержащие указанный текст
// @Tags comments
// @Accept json
// @Produce json
// @Param search body dto.SearchText true "Текст для поиска в комментариях"
// @Success 200 {array} model.Comment "Найденные комментарии"
// @Failure 400 {object} map[string]string "error":"invalid payload"
// @Failure 500 {object} map[string]string "error":"could not get comments"
// @Router /comments/search [post]
func (h *Handler) GetCommentsByTextSearch(c *ginext.Context) {
	var searchText dto.SearchText
	if err := c.BindJSON(&searchText); err != nil {
		zlog.Logger.Error().Msgf("could unmarshal query body to model: %s", err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid payload: " + err.Error(),
		})
		return
	}

	comments, err := h.service.GetCommentsByTextSearch(searchText.Text)
	if err != nil {
		zlog.Logger.Error().Msgf("could not get comments searched by text: %s", err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get comments: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// @Summary Получить все комментарии
// @Description Возвращает полный список всех комментариев без фильтров c пагинацией
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {array} model.Comment "Все комментарии"
// @Failure 400 {object} map[string]string "error":"could not get comments"
// @Router /comments/all [get]
func (h *Handler) GetAllComments(c *ginext.Context) {
	comment, err := h.service.GetAllComments()
	if err != nil {
		zlog.Logger.Error().Msgf("could not get all comments from db: %s", err.Error())
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "could not get comments: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comment)
}

// GetMainPage godoc
// @Summary Get main page
// @Description Serve the main index.html page
// @Tags main
// @Produce html
// @Success 200 {string} string "HTML page"
// @Router / [get]
func (h *Handler) GetMainPage(c *ginext.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
