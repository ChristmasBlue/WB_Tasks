package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wb-go/wbf/ginext"
	"task-3/internal/dto"
	"task-3/internal/model"
	"task-3/internal/repository"
)

type MockCommentService struct {
	mock.Mock
}

func (m *MockCommentService) GetAllComments() ([]*model.Comment, error) {
	args := m.Called()
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockCommentService) GetCommentsById(id int) ([]*model.Comment, error) {
	args := m.Called(id)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockCommentService) GetCommentsPaginated(config dto.CommentsPagination) ([]*model.Comment, error) {
	args := m.Called(config)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockCommentService) GetCommentsByTextSearch(text string) ([]*model.Comment, error) {
	args := m.Called(text)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockCommentService) CreateComment(comment dto.CreateComment) (*dto.CreateComment, error) {
	args := m.Called(comment)
	return args.Get(0).(*dto.CreateComment), args.Error(1)
}

func (m *MockCommentService) DeleteCommentById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestHandler_CreateComment_Success(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	createdAt := time.Date(2025, 9, 20, 0, 0, 0, 0, time.UTC)
	comment := dto.CreateComment{
		ID:        1,
		ParentID:  nil,
		Text:      "Test comment",
		CreatedAt: createdAt,
	}
	expected := &comment

	mockService.On("CreateComment", comment).Return(expected, nil)

	body, _ := json.Marshal(comment)
	req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateComment((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response dto.CreateComment
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expected, &response)
	mockService.AssertExpectations(t)
}

func TestHandler_CreateComment_InvalidPayload(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateComment((*ginext.Context)(c))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "CreateComment")
}

func TestHandler_CreateComment_ServiceError(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	comment := dto.CreateComment{
		Text: "Test comment",
	}

	mockService.On("CreateComment", comment).Return((*dto.CreateComment)(nil), errors.New("service error"))

	body, _ := json.Marshal(comment)
	req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateComment((*ginext.Context)(c))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_CreateComment_InvalidParentID(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	comment := dto.CreateComment{
		Text: "Test comment",
	}

	mockService.On("CreateComment", comment).Return((*dto.CreateComment)(nil), repository.ErrInvalidParenID)

	body, _ := json.Marshal(comment)
	req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.CreateComment((*ginext.Context)(c))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_DeleteCommentById_Success(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	id := "1"

	mockService.On("DeleteCommentById", 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/comments/"+id, nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}

	handler.DeleteCommentById((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_DeleteCommentById_InvalidID(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	id := "invalid"

	req := httptest.NewRequest(http.MethodDelete, "/comments/"+id, nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}

	handler.DeleteCommentById((*ginext.Context)(c))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "DeleteCommentById")
}

func TestHandler_DeleteCommentById_ServiceError(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	id := "1"

	mockService.On("DeleteCommentById", 1).Return(errors.New("service error"))

	req := httptest.NewRequest(http.MethodDelete, "/comments/"+id, nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "id", Value: id}}

	handler.DeleteCommentById((*ginext.Context)(c))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetAllComments_Success(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	expected := []*model.Comment{
		{ID: 1, Text: "Comment 1"},
	}

	mockService.On("GetAllComments").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/comments/all", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetAllComments((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response []*model.Comment
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expected, response)
	mockService.AssertExpectations(t)
}

func TestHandler_GetAllComments_Error(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	mockService.On("GetAllComments").Return(([]*model.Comment)(nil), errors.New("service error"))

	req := httptest.NewRequest(http.MethodGet, "/comments/all", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetAllComments((*ginext.Context)(c))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetCommentsByTextSearch_Success(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	searchText := dto.SearchText{Text: "search"}
	expected := []*model.Comment{
		{ID: 1, Text: "Comment with search"},
	}

	mockService.On("GetCommentsByTextSearch", searchText.Text).Return(expected, nil)

	body, _ := json.Marshal(searchText)
	req := httptest.NewRequest(http.MethodPost, "/comments/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetCommentsByTextSearch((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response []*model.Comment
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expected, response)
	mockService.AssertExpectations(t)
}

func TestHandler_GetCommentsByTextSearch_InvalidPayload(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	req := httptest.NewRequest(http.MethodPost, "/comments/search", bytes.NewBufferString("invalid"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetCommentsByTextSearch((*ginext.Context)(c))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertNotCalled(t, "GetCommentsByTextSearch")
}

func TestHandler_GetCommentsByTextSearch_Error(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	searchText := dto.SearchText{Text: "search"}

	mockService.On("GetCommentsByTextSearch", searchText.Text).Return(([]*model.Comment)(nil), errors.New("service error"))

	body, _ := json.Marshal(searchText)
	req := httptest.NewRequest(http.MethodPost, "/comments/search", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetCommentsByTextSearch((*ginext.Context)(c))

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestHandler_GetComments_Success(t *testing.T) {
	mockService := &MockCommentService{}
	handler := New(mockService)

	expected := []*model.Comment{
		{ID: 1, Text: "Comment 1"},
	}

	mockService.On("GetCommentsById", 1).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/comments?parent=1", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req

	handler.GetComments((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response []*model.Comment
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, expected, response)
	mockService.AssertExpectations(t)
}
