package test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vkr-mtuci/jira-service/config"
	"github.com/vkr-mtuci/jira-service/internal/adapter"
	"github.com/vkr-mtuci/jira-service/internal/handler"
	"github.com/vkr-mtuci/jira-service/internal/service"
)

// MockJiraService - мок-сервис для HTTP-тестирования
type MockJiraService struct {
	mock.Mock
}

// GetIssueDetails - мок-метод для тестирования `JiraService`
func (m *MockJiraService) GetIssueDetails(issueID string) (*adapter.IssueResponse, error) {
	args := m.Called(issueID)
	if args.Get(0) != nil {
		return args.Get(0).(*adapter.IssueResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// TestGetIssueHandler_Success - тест успешного получения задачи через HTTP
func TestGetIssueHandler_Success(t *testing.T) {
	mockService := new(MockJiraService)
	handler := handler.NewJiraHandler(mockService)
	app := fiber.New()
	app.Get("/issue/:id", handler.GetIssue)

	expectedIssue := &adapter.IssueResponse{
		Key: "TEST-123",
		Fields: adapter.IssueFields{
			Summary: "Тестовая задача",
			Status: struct {
				Name string `json:"name"`
			}{Name: "В работе"},
		},
	}

	mockService.On("GetIssueDetails", "TEST-123").Return(expectedIssue, nil)

	req := httptest.NewRequest(http.MethodGet, "/issue/TEST-123", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// TestGetIssueHandler_NotFound - тест, когда задача не найдена
func TestGetIssueHandler_NotFound(t *testing.T) {
	mockService := new(MockJiraService)
	handler := handler.NewJiraHandler(mockService)
	app := fiber.New()
	app.Get("/issue/:id", handler.GetIssue)

	mockService.On("GetIssueDetails", "TEST-404").Return((*adapter.IssueResponse)(nil), errors.New("Задача не найдена"))

	req := httptest.NewRequest(http.MethodGet, "/issue/TEST-404", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// TestGetMultipleIssuesHandler_Success - тест запроса нескольких задач
func TestGetMultipleIssuesHandler_Success(t *testing.T) {
	mockService := new(MockJiraService)
	handler := handler.NewJiraHandler(mockService)
	app := fiber.New()
	app.Get("/issues", handler.GetMultipleIssues)

	expectedIssue := &adapter.IssueResponse{
		Key: "TEST-123",
		Fields: adapter.IssueFields{
			Summary: "Тестовая задача",
			Status: struct {
				Name string `json:"name"`
			}{Name: "В работе"},
		},
	}

	mockService.On("GetIssueDetails", "TEST-123").Return(expectedIssue, nil)
	mockService.On("GetIssueDetails", "TEST-124").Return((*adapter.IssueResponse)(nil), errors.New("Задача не найдена"))

	req := httptest.NewRequest(http.MethodGet, "/issues?issueIDs=TEST-123,TEST-124", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetIssue_JiraDown(t *testing.T) {
	// 🛠 Создаём моковый сервер, который возвращает 500
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	// 🛠 Конфигурация с моковым URL
	cfg := &config.Config{
		JiraBaseURL:  mockServer.URL,
		JiraAPIUrl:   "/rest/api/latest/",
		JiraAPIToken: "dummy-token",
	}

	client := adapter.NewJiraClient(cfg)
	jiraService := service.NewJiraService(client)
	jiraHandler := handler.NewJiraHandler(jiraService)

	// 🛠 Запрос через Fiber
	app := fiber.New()
	app.Get("/issue/:id", jiraHandler.GetIssue)

	req := httptest.NewRequest(http.MethodGet, "/issue/TEST-123", nil)
	resp, _ := app.Test(req)

	// ✅ Проверяем статус 500
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetMultipleIssues_PartialSuccess(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/rest/api/latest/issue/TEST-FAIL" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key": "TEST-123", "fields": {"summary": "Test issue"}}`))
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		JiraBaseURL:  mockServer.URL,
		JiraAPIUrl:   "/rest/api/latest/",
		JiraAPIToken: "dummy-token",
	}
	client := adapter.NewJiraClient(cfg)
	jiraService := service.NewJiraService(client)
	jiraHandler := handler.NewJiraHandler(jiraService)

	app := fiber.New()
	app.Get("/issues", jiraHandler.GetMultipleIssues)

	req := httptest.NewRequest(http.MethodGet, "/issues?issueIDs=TEST-123,TEST-FAIL", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
