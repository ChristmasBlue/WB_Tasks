package service

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"task-2/internal/dto"
	"task-2/internal/model"
)

// MockStorage is a mock implementation of the Storage interface
type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) CreateShortUrl(url model.Url) (*model.Url, error) {
	args := m.Called(url)
	return args.Get(0).(*model.Url), args.Error(1)
}

func (m *MockStorage) CreateRedirectInfo(redirectInfo model.RedirectInfo) error {
	args := m.Called(redirectInfo)
	return args.Error(0)
}

func (m *MockStorage) GetUrlByShort(short string, redirectInfo model.RedirectInfo) (*model.Url, error) {
	args := m.Called(short, redirectInfo)
	return args.Get(0).(*model.Url), args.Error(1)
}

func (m *MockStorage) GetAnalytics(short_url string) ([]dto.RedirectInfo, error) {
	args := m.Called(short_url)
	return args.Get(0).([]dto.RedirectInfo), args.Error(1)
}

func (m *MockStorage) AggregateByUserAgent() ([]dto.UserAgentDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.UserAgentDTO), args.Error(1)
}

func (m *MockStorage) AggregateByDate() ([]dto.DateDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.DateDTO), args.Error(1)
}

func (m *MockStorage) AggregateByMonth() ([]dto.MonthDTO, error) {
	args := m.Called()
	return args.Get(0).([]dto.MonthDTO), args.Error(1)
}

// MockCache is a mock implementation of the Cache interface
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Set(key string, value interface{}) error {
	args := m.Called(key, value)
	return args.Error(0)
}

func TestService_CreateShortUrl_CacheHit(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	url := model.Url{Url: "https://example.com"}
	shortUrl := "abc123"

	mockCache.On("Get", url.Url).Return(shortUrl, nil)

	result, err := service.CreateShortUrl(url)

	assert.NoError(t, err)
	assert.Equal(t, url.Url, result.Url)
	assert.Equal(t, shortUrl, result.ShortUrl)
	mockCache.AssertExpectations(t)
	mockStorage.AssertNotCalled(t, "CreateShortUrl", mock.Anything)
}

func TestService_CreateShortUrl_CacheMiss(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	url := model.Url{Url: "https://example.com"}
	shortUrl := "def456"
	createdUrl := &model.Url{Url: url.Url, ShortUrl: shortUrl}

	mockCache.On("Get", url.Url).Return("", redis.Nil)
	mockStorage.On("CreateShortUrl", mock.AnythingOfType("model.Url")).Return(createdUrl, nil).Once()

	result, err := service.CreateShortUrl(url)

	assert.NoError(t, err)
	assert.Equal(t, createdUrl, result)
	mockCache.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestService_GetUrlByShort_CacheHit(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	shortUrl := "abc123"
	originalUrl := "https://example.com"
	redirectInfo := model.RedirectInfo{ShortUrl: shortUrl}

	mockCache.On("Get", shortUrl).Return(originalUrl, nil)
	mockStorage.On("CreateRedirectInfo", redirectInfo).Return(nil)

	result, err := service.GetUrlByShort(shortUrl, redirectInfo)

	assert.NoError(t, err)
	assert.Equal(t, shortUrl, result.ShortUrl)
	assert.Equal(t, originalUrl, result.Url)
	mockCache.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestService_GetUrlByShort_CacheMiss(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	shortUrl := "abc123"
	originalUrl := "https://example.com"
	redirectInfo := model.RedirectInfo{ShortUrl: shortUrl}
	urlInfo := &model.Url{ShortUrl: shortUrl, Url: originalUrl}

	mockCache.On("Get", shortUrl).Return("", redis.Nil)
	mockStorage.On("GetUrlByShort", shortUrl, redirectInfo).Return(urlInfo, nil)

	result, err := service.GetUrlByShort(shortUrl, redirectInfo)

	assert.NoError(t, err)
	assert.Equal(t, urlInfo, result)
	mockCache.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestService_GetAnalytics(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	shortUrl := "abc123"
	analytics := []dto.RedirectInfo{
		{Url: "https://example.com", ShortUrl: shortUrl, RedirectCount: 10},
	}

	mockStorage.On("GetAnalytics", shortUrl).Return(analytics, nil)
	mockCache.On("Set", "https://example.com", shortUrl).Return(nil)

	result, err := service.GetAnalytics(shortUrl)

	assert.NoError(t, err)
	assert.Equal(t, analytics, result)
	mockStorage.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestService_AggregateByUserAgent(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	expected := []dto.UserAgentDTO{
		{ShortUrl: "abc123", UserAgent: []string{"Mozilla/5.0"}, RedirectCount: 5},
	}

	mockStorage.On("AggregateByUserAgent").Return(expected, nil)

	result, err := service.AggregateByUserAgent()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}

func TestService_AggregateByDate(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	expected := []dto.DateDTO{
		{Day: 1, Month: 1, Year: 2023, UrlInfo: []dto.UrlInfo{{ShortUrl: "abc123", Time: "10:00"}}, RedirectCount: 10},
	}

	mockStorage.On("AggregateByDate").Return(expected, nil)

	result, err := service.AggregateByDate()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}

func TestService_AggregateByMonth(t *testing.T) {
	mockStorage := new(MockStorage)
	mockCache := new(MockCache)
	service := New(mockStorage, mockCache)

	expected := []dto.MonthDTO{
		{Month: 1, Year: 2023, UrlInfo: []struct {
			ShortUrl string "json:\"short_url\""
			Time     string "json:\"time\""
		}{{ShortUrl: "abc123", Time: "10:00"}}, RedirectCount: 100},
	}

	mockStorage.On("AggregateByMonth").Return(expected, nil)

	result, err := service.AggregateByMonth()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
	mockStorage.AssertExpectations(t)
}
