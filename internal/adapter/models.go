package adapter

import (
	"context"

	"github.com/go-resty/resty/v2"
)

// IssueFields - основные поля задачи Jira
type IssueFields struct {
	Summary     string `json:"summary"`     // Краткое описание
	Description string `json:"description"` // Полное описание
	Status      struct {
		Name string `json:"name"` // Статус задачи
	} `json:"status"`
	Assignee struct {
		DisplayName string `json:"displayName"` // Исполнитель
	} `json:"assignee"`
	Reporter struct {
		DisplayName string `json:"displayName"` // Автор задачи
	} `json:"reporter"`
	Priority struct {
		Name string `json:"name"` // Приоритет
	} `json:"priority"`
	IssueType struct {
		Name string `json:"name"` // Тип задачи
	} `json:"issuetype"`
	Project struct {
		Name string `json:"name"` // Название проекта
	} `json:"project"`
	Created string `json:"created"` // Дата создания
	Updated string `json:"updated"` // Дата обновления
}

// IssueResponse - структура ответа Jira
type IssueResponse struct {
	Key    string      `json:"key"`    // Идентификатор задачи
	Fields IssueFields `json:"fields"` // Данные задачи
}

// JiraErrorResponse - обработка ошибок от Jira API
type JiraErrorResponse struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
}

// 🔹 Интерфейс для мока Jira-клиента
type JiraClientInterface interface {
	GetIssue(ctx context.Context, issueID string) (*IssueResponse, error)
}

// JiraClient - клиент для работы с Jira API
type JiraClient struct {
	client  *resty.Client
	baseURL string
	apiURL  string
}
