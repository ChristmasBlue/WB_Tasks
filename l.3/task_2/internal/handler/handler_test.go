package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wb-go/wbf/ginext"
	"task-2/internal/dto"
	"task-2/internal/model"
)

// MockShortnerService is a mock implementation of the ShortnerService interface
type MockShortnerService struct {
	mock.Mock
}

func (m *MockShortnerService) CreateShortUrl(url model.Url) (*model.Url, error) {
	args := m.Called(url)
	return args.Get(0).(*model.Url), args.Error(1)
}

func (m *MockShortnerService) GetUrlByShort(short_url string, redirectInfo model.RedirectInfo) (*model.Url, error) {
	args := m.Called(short_url, redirectInfo)
	return args.Get(0).(*model.Url), args.Error(1)
}

func (m *MockShortnerService) GetAnalytics(short_url string) ([]dto.RedirectInfo, error) {
	args := m.Called(short_url)
	return args.Get(0).([]dto.RedirectInfo), args.Error(1)
}

func (m *MockShortnerService) AggregateByUserAgent() ([]dto.UserAgentDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.UserAgentDTO), args.Error(1)
}

func (m *MockShortnerService) AggregateByDate() ([]dto.DateDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.DateDTO), args.Error(1)
}

func (m *MockShortnerService) AggregateByMonth() ([]dto.MonthDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.MonthDTO), args.Error(1)
}

func TestHandler_CreateShortUrl_Success(t *testing.T) {
	mockService := new(MockShortnerService)
	handler := New(mockService)

	url := model.Url{Url: "https://example.com"}
	shortUrl := &model.Url{Url: url.Url, ShortUrl: "abc123"}

	mockService.On("CreateShortUrl", url).Return(shortUrl, nil)

	reqBody := `{"url": "https://example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.CreateShortUrl((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response model.Url
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, shortUrl, &response)
	mockService.AssertExpectations(t)
}

func TestHandler_CreateShortUrl_InvalidJSON(t *testing.T) {
	mockService := new(MockShortnerService)
	handler := New(mockService)

	req := httptest.NewRequest(http.MethodPost, "/shorten", strings.NewReader("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.CreateShortUrl((*ginext.Context)(c))

	assert.Equal(t, http.StatusBadRequest, w.Code)
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid request body", response["error"])
}

func TestHandler_RedirectByShortUrl_Success(t *testing.T) {
	mockService := new(MockShortnerService)
	handler := New(mockService)

	shortUrl := "abc123"
	originalUrl := "https://example.com"
	redirectInfo := model.RedirectInfo{ShortUrl: shortUrl, UserAgent: "test-agent"}
	urlInfo := &model.Url{ShortUrl: shortUrl, Url: originalUrl}

	mockService.On("GetUrlByShort", shortUrl, redirectInfo).Return(urlInfo, nil)

	req := httptest.NewRequest(http.MethodGet, "/s/"+shortUrl, nil)
	req.Header.Set("User-Agent", "test-agent")
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "short_url", Value: shortUrl}}
	handler.RedirectByShortUrl((*ginext.Context)(c))

	assert.Equal(t, http.StatusMovedPermanently, w.Code)
	assert.Equal(t, originalUrl, w.Header().Get("Location"))
	mockService.AssertExpectations(t)
}

func TestHandler_GetAnalytics_Success(t *testing.T) {
	mockService := new(MockShortnerService)
	handler := New(mockService)

	shortUrl := "abc123"
	analytics := []dto.RedirectInfo{
		{Url: "https://example.com", ShortUrl: shortUrl, RedirectCount: 5, RequestTime: []string{"10:00"}, UserAgent: []string{"Mozilla/5.0"}},
	}

	mockService.On("GetAnalytics", shortUrl).Return(analytics, nil)

	req := httptest.NewRequest(http.MethodGet, "/analytics/"+shortUrl, nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{{Key: "short_url", Value: shortUrl}}
	handler.GetAnalytics((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response []dto.RedirectInfo
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, analytics, response)
	mockService.AssertExpectations(t)
}

func TestHandler_AggregateByUserAgent_Success(t *testing.T) {
	mockService := new(MockShortnerService)
	handler := New(mockService)

	expected := []dto.UserAgentDTO{
		{ShortUrl: "abc123", UserAgent: []string{"Mozilla/5.0"}, RedirectCount: 5},
	}

	mockService.On("AggregateByUserAgent").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/analytics/user_agent", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.AggregateByUserAgent((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response []dto.UserAgentDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expected, response)
	mockService.AssertExpectations(t)
}

func TestHandler_AggregateByDate_Success(t *testing.T) {
	mockService := new(MockShortnerService)
	handler := New(mockService)

	expected := []dto.DateDTO{
		{Day: 1, Month: 1, Year: 2023, UrlInfo: []dto.UrlInfo{{ShortUrl: "abc123", Time: "10:00"}}, RedirectCount: 10},
	}

	mockService.On("AggregateByDate").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/analytics/date", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.AggregateByDate((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response []dto.DateDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expected, response)
	mockService.AssertExpectations(t)
}

func TestHandler_AggregateByMonth_Success(t *testing.T) {
	mockService := new(MockShortnerService)
	handler := New(mockService)

	expected := []dto.MonthDTO{
		{Month: 1, Year: 2023, UrlInfo: []struct {
			ShortUrl string "json:\"short_url\""
			Time     string "json:\"time\""
		}{{ShortUrl: "abc123", Time: "10:00"}}, RedirectCount: 100},
	}

	mockService.On("AggregateByMonth").Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/analytics/month", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	handler.AggregateByMonth((*ginext.Context)(c))

	assert.Equal(t, http.StatusOK, w.Code)
	var response []dto.MonthDTO
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, expected, response)
	mockService.AssertExpectations(t)
}
