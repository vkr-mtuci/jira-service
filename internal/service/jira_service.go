package service

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/vkr-mtuci/jira-service/internal/adapter"
)

// 🔹 Интерфейс для JiraService (чтобы можно было мокать)
type JiraServiceInterface interface {
	GetIssueDetails(issueID string) (*adapter.IssueResponse, error)
}

// JiraService - сервис для работы с задачами Jira
type JiraService struct {
	client adapter.JiraClientInterface
}

// ✅ Теперь JiraService принимает не `*adapter.JiraClient`, а `JiraClientInterface`
func NewJiraService(client adapter.JiraClientInterface) *JiraService {
	return &JiraService{client: client}
}

// GetIssueDetails получает информацию о задаче в Jira
func (s *JiraService) GetIssueDetails(issueID string) (*adapter.IssueResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	issue, err := s.client.GetIssue(ctx, issueID)
	if err != nil {
		log.Error().Err(err).Msgf("❌ Ошибка получения задачи %s", issueID)
		return nil, err
	}

	log.Info().Msgf("✅ Успешно получена информация о задаче %s", issueID)
	return issue, nil
}
