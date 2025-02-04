package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
	"github.com/vkr-mtuci/jira-service/config"
)

// NewJiraClient - инициализация клиента
func NewJiraClient(cfg *config.Config) *JiraClient {
	client := resty.New().
		SetBaseURL(cfg.JiraBaseURL).
		SetTimeout(10*time.Second).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetAuthToken(cfg.JiraAPIToken) // ✅ Bearer Token

	log.Info().Msg("🔗 Подключение к Jira API: " + cfg.JiraBaseURL)

	return &JiraClient{
		client:  client,
		baseURL: cfg.JiraBaseURL,
		apiURL:  cfg.JiraAPIUrl,
	}
}

// GetIssue - получение информации о задаче Jira
func (j *JiraClient) GetIssue(ctx context.Context, issueID string) (*IssueResponse, error) {
	if issueID == "" {
		return nil, fmt.Errorf("❌ issueID не может быть пустым")
	}

	url := fmt.Sprintf("%s%sissue/%s", j.baseURL, j.apiURL, issueID)

	log.Debug().Msgf("📡 Запрос в Jira: issueID=%s, URL=%s", issueID, url)

	// Запрос с учетом контекста (таймаут, отмена)
	resp, err := j.client.R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		log.Error().Err(err).Msg("❌ Ошибка запроса в Jira")
		return nil, err
	}

	// Обрабатываем HTTP-статусы
	switch resp.StatusCode() {
	case http.StatusOK:
		var issue IssueResponse
		if err := json.Unmarshal(resp.Body(), &issue); err != nil {
			log.Error().Err(err).Msg("❌ Ошибка парсинга JSON из Jira")
			return nil, err
		}
		log.Info().Msgf("✅ Получена задача %s", issueID)
		return &issue, nil
	case http.StatusUnauthorized:
		return nil, fmt.Errorf("⚠️ 401 Unauthorized: проверьте API-токен")
	case http.StatusForbidden:
		return nil, fmt.Errorf("⚠️ 403 Forbidden: нет доступа к задаче %s", issueID)
	case http.StatusNotFound:
		return nil, fmt.Errorf("⚠️ 404 Not Found: задача %s не найдена", issueID)
	default:
		var jiraErr JiraErrorResponse
		_ = json.Unmarshal(resp.Body(), &jiraErr)
		return nil, fmt.Errorf("⚠️ Ошибка от Jira: %s", jiraErr.ErrorMessages)
	}
}
