package handler

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"task-3/internal/dto"
	"task-3/internal/repository"
)

func parserQueryParameters(params url.Values) (*dto.CommentsPagination, error) {
	var id, page, limit int
	var err error
	for param, value := range params {
		if param == "parent" && len(value) != 0 {
			id, err = strconv.Atoi(value[0])
			if err != nil {
				return nil, fmt.Errorf("invalid limit was provided: %s", err.Error())
			}
		}

		if param == "page" && len(value) != 0 {
			page, _ = strconv.Atoi(value[0])
		}

		if param == "limit" && len(value) != 0 {
			limit, _ = strconv.Atoi(value[0])
		}
	}

	var config dto.CommentsPagination
	config.ParentID = id
	config.Page = page
	config.Limit = limit

	return &config, nil
}

func (h *Handler) getCommentsById(id int, c *ginext.Context) {
	comments, err := h.service.GetCommentsById(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotSuchComment) {
			zlog.Logger.Error().Msgf("invalid id: %s", err.Error())
			c.JSON(http.StatusBadRequest, ginext.H{
				"error": "invalid id: " + err.Error(),
			})
			return
		}

		zlog.Logger.Error().Msgf("could not get comment by id: %s", err.Error())
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": "could not get comment: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, comments)
}
