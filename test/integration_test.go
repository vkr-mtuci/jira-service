package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vkr-mtuci/jira-service/internal/adapter"
	"github.com/vkr-mtuci/jira-service/internal/handler"
	"github.com/vkr-mtuci/jira-service/internal/service"
)

func TestFullIntegration(t *testing.T) {
	// Создаем мок-клиент
	mockClient := new(MockJiraClient)

	// Настраиваем мок: при запросе "TEST-987" возвращаем фиктивный ответ
	mockResponse := &adapter.IssueResponse{
		Key: "TEST-987",
		Fields: adapter.IssueFields{
			Summary: "Mocked Jira Issue",
		},
	}
	mockClient.On("GetIssue", mock.Anything, "TEST-987").Return(mockResponse, nil)

	// Создаем сервис и обработчик с мок-клиентом
	jiraService := service.NewJiraService(mockClient)
	jiraHandler := handler.NewJiraHandler(jiraService)

	// Создаем Fiber-приложение
	app := fiber.New()
	app.Get("/issue/:id", jiraHandler.GetIssue)

	// Запускаем тестовый запрос
	req := httptest.NewRequest(http.MethodGet, "/issue/TEST-987", nil)
	resp, _ := app.Test(req)

	// Проверяем, что вернулся 200 OK
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Проверяем, что мок-метод был вызван
	mockClient.AssertExpectations(t)
}
