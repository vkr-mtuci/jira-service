package test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vkr-mtuci/jira-service/internal/adapter"
	"github.com/vkr-mtuci/jira-service/internal/service"
)

// MockJiraClient - мок-реализация клиента Jira
type MockJiraClient struct {
	mock.Mock
}

// GetIssue - мок-метод для тестирования `JiraClient.GetIssue`
func (m *MockJiraClient) GetIssue(ctx context.Context, issueID string) (*adapter.IssueResponse, error) {
	args := m.Called(ctx, issueID)
	if args.Get(0) != nil {
		return args.Get(0).(*adapter.IssueResponse), args.Error(1)
	}
	return nil, args.Error(1)
}

// TestGetIssueDetails_Success - успешный тест получения задачи
func TestGetIssueDetails_Success(t *testing.T) {
	mockClient := new(MockJiraClient)
	service := service.NewJiraService(mockClient)

	expectedIssue := &adapter.IssueResponse{
		Key: "TEST-123",
		Fields: adapter.IssueFields{
			Summary: "Тестовая задача",
			Status: struct {
				Name string `json:"name"`
			}{Name: "В работе"},
		},
	}

	mockClient.On("GetIssue", mock.Anything, "TEST-123").Return(expectedIssue, nil)

	issue, err := service.GetIssueDetails("TEST-123")

	assert.NoError(t, err)
	assert.Equal(t, expectedIssue, issue)

	mockClient.AssertExpectations(t)
}

// TestGetIssueDetails_Error - тест на ошибку (задача не найдена)
func TestGetIssueDetails_Error(t *testing.T) {
	mockClient := new(MockJiraClient)
	service := service.NewJiraService(mockClient)

	mockClient.On("GetIssue", mock.Anything, "TEST-404").Return((*adapter.IssueResponse)(nil), errors.New("Задача не найдена"))

	issue, err := service.GetIssueDetails("TEST-404")

	assert.Error(t, err)
	assert.Nil(t, issue)

	mockClient.AssertExpectations(t)
}
