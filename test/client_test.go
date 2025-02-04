package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkr-mtuci/jira-service/config"
	"github.com/vkr-mtuci/jira-service/internal/adapter"
)

// 📌 Тест: успешное получение задачи
func TestGetIssue_Success(t *testing.T) {
	// 🛠 Моковый сервер
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"key": "TEST-123", "fields": {"summary": "Test issue"}}`))
	}))
	defer mockServer.Close()

	// 🛠 Конфигурация с моковым URL
	cfg := &config.Config{
		JiraBaseURL:  mockServer.URL,
		JiraAPIUrl:   "/rest/api/latest/",
		JiraAPIToken: "dummy-token",
	}

	client := adapter.NewJiraClient(cfg)

	// 📡 Запрос
	ctx := context.TODO()
	issue, err := client.GetIssue(ctx, "TEST-123")

	// ✅ Проверяем
	assert.NoError(t, err)
	assert.NotNil(t, issue)
	assert.Equal(t, "TEST-123", issue.Key)
}

func TestGetIssue_ErrorCases(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/rest/api/latest/issue/TEST-404":
			w.WriteHeader(http.StatusNotFound)
		case "/rest/api/latest/issue/TEST-401":
			w.WriteHeader(http.StatusUnauthorized)
		case "/rest/api/latest/issue/TEST-403":
			w.WriteHeader(http.StatusForbidden)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer mockServer.Close()

	cfg := &config.Config{
		JiraBaseURL:  mockServer.URL,
		JiraAPIUrl:   "/rest/api/latest/",
		JiraAPIToken: "dummy-token",
	}
	client := adapter.NewJiraClient(cfg)

	ctx := context.TODO()

	_, err := client.GetIssue(ctx, "TEST-404")
	assert.Error(t, err, "Должна быть ошибка 404")

	_, err = client.GetIssue(ctx, "TEST-401")
	assert.Error(t, err, "Должна быть ошибка 401")

	_, err = client.GetIssue(ctx, "TEST-403")
	assert.Error(t, err, "Должна быть ошибка 403")
}
