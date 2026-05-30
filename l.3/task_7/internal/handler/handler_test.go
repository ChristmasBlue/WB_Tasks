package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"task-7/internal/dto"
	"task-7/internal/model"

	"github.com/gin-gonic/gin"
)

type mockService struct {
	createUserFunc          func(ctx context.Context, createUser dto.CreateUser) (*model.User, error)
	createItemFunc          func(ctx context.Context, createItem dto.CreateItem) (*model.Item, error)
	getAllItemsFunc         func(ctx context.Context) ([]model.Item, error)
	getUsersWithChangesFunc func(ctx context.Context) ([]model.UserHistory, error)
	updateItemFunc          func(ctx context.Context, updateItem *dto.UpdateItem) error
	deleteItemFunc          func(ctx context.Context, id int) error
	getUserRoleFunc         func(ctx context.Context, id int) (string, error)
}

func (m *mockService) CreateUser(ctx context.Context, createUser dto.CreateUser) (*model.User, error) {
	return m.createUserFunc(ctx, createUser)
}

func (m *mockService) CreateItem(ctx context.Context, createItem dto.CreateItem) (*model.Item, error) {
	return m.createItemFunc(ctx, createItem)
}

func (m *mockService) GetAllItems(ctx context.Context) ([]model.Item, error) {
	return m.getAllItemsFunc(ctx)
}

func (m *mockService) GetUsersWithChanges(ctx context.Context) ([]model.UserHistory, error) {
	return m.getUsersWithChangesFunc(ctx)
}

func (m *mockService) UpdateItem(ctx context.Context, updateItem *dto.UpdateItem) error {
	return m.updateItemFunc(ctx, updateItem)
}

func (m *mockService) DeleteItem(ctx context.Context, id int) error {
	return m.deleteItemFunc(ctx, id)
}

func (m *mockService) GetUserRole(ctx context.Context, id int) (string, error) {
	return m.getUserRoleFunc(ctx, id)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}

func TestHandler_CreateUser(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		mockFunc         func(ctx context.Context, createUser dto.CreateUser) (*model.User, error)
		expectedStatus   int
		checkToken       bool
		expectedResponse map[string]interface{}
	}{
		{
			name: "success",
			body: `{"name":"test","role":"admin"}`,
			mockFunc: func(ctx context.Context, createUser dto.CreateUser) (*model.User, error) {
				return &model.User{ID: 1, Name: "test", Role: "admin", CreatedAt: time.Now()}, nil
			},
			expectedStatus: http.StatusOK,
			checkToken:     true,
			expectedResponse: map[string]interface{}{
				"id":   float64(1),
				"name": "test",
				"role": "admin",
			},
		},
		{
			name:           "bind error",
			body:           `invalid json`,
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid character 'i' looking for beginning of value",
			},
		},
		{
			name:           "validation error",
			body:           `{"name":"","role":"admin"}`,
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "Name is required",
			},
		},
		{
			name: "service error",
			body: `{"name":"test","role":"admin"}`,
			mockFunc: func(ctx context.Context, createUser dto.CreateUser) (*model.User, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SECRET", "testsecret")
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			mock := &mockService{
				createUserFunc: tt.mockFunc,
			}
			h := New(context.Background(), mock)
			h.CreateUser(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			for k, v := range tt.expectedResponse {
				if resp[k] != v {
					t.Errorf("expected %v for %s, got %v", v, k, resp[k])
				}
			}

			if tt.checkToken {
				if _, ok := resp["token"]; !ok {
					t.Error("expected token in response")
				}
			}
		})
	}
}

func TestHandler_CreateItem(t *testing.T) {
	tests := []struct {
		name             string
		body             string
		mockFunc         func(ctx context.Context, createItem dto.CreateItem) (*model.Item, error)
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name: "success",
			body: `{"name":"item","count":10}`,
			mockFunc: func(ctx context.Context, createItem dto.CreateItem) (*model.Item, error) {
				return &model.Item{ID: 1, Name: "item", Count: 10, CreatedAt: time.Now()}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"id":    float64(1),
				"name":  "item",
				"count": float64(10),
			},
		},
		{
			name:           "bind error",
			body:           `invalid json`,
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid character 'i' looking for beginning of value",
			},
		},
		{
			name:           "validation error",
			body:           `{"name":"","count":10}`,
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "Name is required",
			},
		},
		{
			name: "service error",
			body: `{"name":"item","count":10}`,
			mockFunc: func(ctx context.Context, createItem dto.CreateItem) (*model.Item, error) {
				return nil, errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")

			mock := &mockService{
				createItemFunc: tt.mockFunc,
			}
			h := New(context.Background(), mock)
			h.CreateItem(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			for k, v := range tt.expectedResponse {
				if resp[k] != v {
					t.Errorf("expected %v for %s, got %v", v, k, resp[k])
				}
			}
		})
	}
}

func TestHandler_GetAllItems(t *testing.T) {
	tests := []struct {
		name             string
		mockFunc         func(ctx context.Context) ([]model.Item, error)
		expectedStatus   int
		expectedResponse []map[string]interface{}
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context) ([]model.Item, error) {
				return []model.Item{{ID: 1, Name: "item", Count: 10, CreatedAt: time.Now()}}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []map[string]interface{}{
				{
					"id":    float64(1),
					"name":  "item",
					"count": float64(10),
				},
			},
		},
		{
			name: "service error",
			mockFunc: func(ctx context.Context) ([]model.Item, error) {
				return nil, errors.New("service error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/items", nil)

			mock := &mockService{
				getAllItemsFunc: tt.mockFunc,
			}
			h := New(context.Background(), mock)
			h.GetAllItems(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedResponse != nil {
				var resp []map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				for i, expected := range tt.expectedResponse {
					for k, v := range expected {
						if resp[i][k] != v {
							t.Errorf("expected %v for %s, got %v", v, k, resp[i][k])
						}
					}
				}
			} else {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp["error"] != "service error" {
					t.Errorf("expected error 'service error', got %v", resp["error"])
				}
			}
		})
	}
}

func TestHandler_GetUsersWithChanges(t *testing.T) {
	tests := []struct {
		name             string
		mockFunc         func(ctx context.Context) ([]model.UserHistory, error)
		expectedStatus   int
		expectedResponse []map[string]interface{}
	}{
		{
			name: "success",
			mockFunc: func(ctx context.Context) ([]model.UserHistory, error) {
				return []model.UserHistory{{User: model.User{ID: 1, Name: "user", Role: "admin", CreatedAt: time.Now()}}}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []map[string]interface{}{
				{
					"id":   float64(1),
					"name": "user",
					"role": "admin",
				},
			},
		},
		{
			name: "service error",
			mockFunc: func(ctx context.Context) ([]model.UserHistory, error) {
				return nil, errors.New("service error")
			},
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/users/history", nil)

			mock := &mockService{
				getUsersWithChangesFunc: tt.mockFunc,
			}
			h := New(context.Background(), mock)
			h.GetUsersWithChanges(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedResponse != nil {
				var resp []map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				for i, expected := range tt.expectedResponse {
					for k, v := range expected {
						if resp[i][k] != v {
							t.Errorf("expected %v for %s, got %v", v, k, resp[i][k])
						}
					}
				}
			} else {
				var resp map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &resp)
				if resp["error"] != "service error" {
					t.Errorf("expected error 'service error', got %v", resp["error"])
				}
			}
		})
	}
}

func TestHandler_UpdateItem(t *testing.T) {
	tests := []struct {
		name             string
		paramID          string
		body             string
		mockFunc         func(ctx context.Context, updateItem *dto.UpdateItem) error
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:    "success",
			paramID: "1",
			body:    `{"user_id":1,"name":"updated"}`,
			mockFunc: func(ctx context.Context, updateItem *dto.UpdateItem) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"status": "successfully updated item",
			},
		},
		{
			name:           "invalid id",
			paramID:        "abc",
			body:           `{"user_id":1}`,
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid item is or was not provided",
			},
		},
		{
			name:           "bind error",
			paramID:        "1",
			body:           `invalid json`,
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid character 'i' looking for beginning of value",
			},
		},
		{
			name:           "validation error",
			paramID:        "1",
			body:           `{"user_id":0}`,
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "UserID is required",
			},
		},
		{
			name:    "service error",
			paramID: "1",
			body:    `{"user_id":1}`,
			mockFunc: func(ctx context.Context, updateItem *dto.UpdateItem) error {
				return errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPut, "/items/"+tt.paramID, strings.NewReader(tt.body))
			c.Request.Header.Set("Content-Type", "application/json")
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}

			mock := &mockService{
				updateItemFunc: tt.mockFunc,
			}
			h := New(context.Background(), mock)
			h.UpdateItem(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			for k, v := range tt.expectedResponse {
				if resp[k] != v {
					t.Errorf("expected %v for %s, got %v", v, k, resp[k])
				}
			}
		})
	}
}

func TestHandler_DeleteItem(t *testing.T) {
	tests := []struct {
		name             string
		paramID          string
		mockFunc         func(ctx context.Context, id int) error
		expectedStatus   int
		expectedResponse map[string]interface{}
	}{
		{
			name:    "success",
			paramID: "1",
			mockFunc: func(ctx context.Context, id int) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]interface{}{
				"status": "successfully deleted item",
			},
		},
		{
			name:           "invalid id",
			paramID:        "abc",
			mockFunc:       nil,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": "invalid item is or was not provided",
			},
		},
		{
			name:    "service error",
			paramID: "1",
			mockFunc: func(ctx context.Context, id int) error {
				return errors.New("service error")
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": "service error",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodDelete, "/items/"+tt.paramID, nil)
			c.Params = gin.Params{{Key: "id", Value: tt.paramID}}

			mock := &mockService{
				deleteItemFunc: tt.mockFunc,
			}
			h := New(context.Background(), mock)
			h.DeleteItem(c)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var resp map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &resp)

			for k, v := range tt.expectedResponse {
				if resp[k] != v {
					t.Errorf("expected %v for %s, got %v", v, k, resp[k])
				}
			}
		})
	}
}
