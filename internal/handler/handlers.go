package handler

import (
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"github.com/vkr-mtuci/jira-service/internal/service"
)

// JiraHandler теперь зависит не от `JiraService`, а от интерфейса `JiraServiceInterface`
type JiraHandler struct {
	service service.JiraServiceInterface
}

// NewJiraHandler создаёт новый обработчик Jira API
func NewJiraHandler(service service.JiraServiceInterface) *JiraHandler {
	return &JiraHandler{service: service}
}

// GetIssue обрабатывает запрос на получение одной задачи
func (h *JiraHandler) GetIssue(c *fiber.Ctx) error {
	issueID := c.Params("id")
	if issueID == "" {
		log.Warn().Msg("⚠️ Не указан issueID")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Необходимо указать issueID",
		})
	}

	issue, err := h.service.GetIssueDetails(issueID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Ошибка получения задачи",
		})
	}

	return c.JSON(issue)
}

// GetMultipleIssues обрабатывает запрос нескольких задач
func (h *JiraHandler) GetMultipleIssues(c *fiber.Ctx) error {
	issueIDsQuery := c.Query("issueIDs")
	if issueIDsQuery == "" {
		log.Warn().Msg("⚠️ issueIDs не указаны")
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Необходимо указать хотя бы один issueID",
		})
	}

	// Разбиваем строку issueIDs в массив
	issueIDs := strings.Split(issueIDsQuery, ",")
	log.Info().Msgf("📡 Запрашиваем задачи: %v", issueIDs)

	// Создаем массив для хранения результатов
	var issues []interface{}
	for _, id := range issueIDs {
		issue, err := h.service.GetIssueDetails(strings.TrimSpace(id))
		if err != nil {
			log.Error().Err(err).Msgf("❌ Ошибка при получении задачи %s", id)
			continue
		}
		issues = append(issues, issue)
	}

	return c.JSON(fiber.Map{
		"issues": issues,
	})
}
