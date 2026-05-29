package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-3/internal/dto"
	"task-3/internal/model"
)

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetCommentsById(id int) ([]*model.Comment, error) {
	args := m.Called(id)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockStorage) GetCommentsPaginated(config dto.CommentsPagination) ([]*model.Comment, error) {
	args := m.Called(config)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockStorage) GetAllComments() ([]*model.Comment, error) {
	args := m.Called()
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockStorage) GetCommentsByTextSearch(text string) ([]*model.Comment, error) {
	args := m.Called(text)
	return args.Get(0).([]*model.Comment), args.Error(1)
}

func (m *MockStorage) CreateComment(comment dto.CreateComment) (*dto.CreateComment, error) {
	args := m.Called(comment)
	return args.Get(0).(*dto.CreateComment), args.Error(1)
}

func (m *MockStorage) DeleteCommentById(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)
	assert.NotNil(t, service)
	assert.Equal(t, mockStorage, service.storage)
}

func TestService_CreateComment(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	comment := dto.CreateComment{
		ID:        1,
		ParentID:  nil,
		Text:      "Test comment",
		CreatedAt: time.Now(),
	}
	expected := &comment

	mockStorage.On("CreateComment", comment).Return(expected, nil)

	result, err := service.CreateComment(comment)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}

func TestService_DeleteCommentById(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	id := 1

	mockStorage.On("DeleteCommentById", id).Return(nil)

	err := service.DeleteCommentById(id)

	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestService_GetAllComments(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	expected := []*model.Comment{
		{ID: 1, Text: "Comment 1"},
	}

	mockStorage.On("GetAllComments").Return(expected, nil)

	result, err := service.GetAllComments()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetCommentsById(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	id := 1
	expected := []*model.Comment{
		{ID: 1, Text: "Comment 1"},
	}

	mockStorage.On("GetCommentsById", id).Return(expected, nil)

	result, err := service.GetCommentsById(id)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetCommentsPaginated(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	config := dto.CommentsPagination{
		ParentID: 0,
		Page:     1,
		Limit:    10,
	}
	expected := []*model.Comment{
		{ID: 1, Text: "Comment 1"},
	}

	mockStorage.On("GetCommentsPaginated", config).Return(expected, nil)

	result, err := service.GetCommentsPaginated(config)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetCommentsByTextSearch(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	text := "search"
	expected := []*model.Comment{
		{ID: 1, Text: "Comment with search"},
	}

	mockStorage.On("GetCommentsByTextSearch", text).Return(expected, nil)

	result, err := service.GetCommentsByTextSearch(text)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}

// Test error cases
func TestService_CreateComment_Error(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	comment := dto.CreateComment{
		Text: "Test comment",
	}

	mockStorage.On("CreateComment", comment).Return((*dto.CreateComment)(nil), errors.New("storage error"))

	result, err := service.CreateComment(comment)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
}

func TestService_DeleteCommentById_Error(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	id := 1

	mockStorage.On("DeleteCommentById", id).Return(errors.New("storage error"))

	err := service.DeleteCommentById(id)

	assert.Error(t, err)
	mockStorage.AssertExpectations(t)
}

func TestService_GetAllComments_Error(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	mockStorage.On("GetAllComments").Return(([]*model.Comment)(nil), errors.New("storage error"))

	result, err := service.GetAllComments()

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetCommentsById_Error(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	id := 1

	mockStorage.On("GetCommentsById", id).Return(([]*model.Comment)(nil), errors.New("storage error"))

	result, err := service.GetCommentsById(id)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetCommentsPaginated_Error(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	config := dto.CommentsPagination{}

	mockStorage.On("GetCommentsPaginated", config).Return(([]*model.Comment)(nil), errors.New("storage error"))

	result, err := service.GetCommentsPaginated(config)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
}

func TestService_GetCommentsByTextSearch_Error(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	text := "search"

	mockStorage.On("GetCommentsByTextSearch", text).Return(([]*model.Comment)(nil), errors.New("storage error"))

	result, err := service.GetCommentsByTextSearch(text)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockStorage.AssertExpectations(t)
}
