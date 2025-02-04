package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"

	"github.com/vkr-mtuci/jira-service/config"
	"github.com/vkr-mtuci/jira-service/internal/adapter"
	"github.com/vkr-mtuci/jira-service/internal/handler"
	"github.com/vkr-mtuci/jira-service/internal/service"
)

func main() {
	// Настроим zerolog: логи будут выводиться в удобочитаемом формате
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	logger := zerolog.New(output).With().Timestamp().Logger() // ✅ Создаем объект логгера

	// Загружаем конфигурацию
	cfg := config.LoadConfig()

	// Вывод информации о запуске сервиса
	logger.Info().Msg("📢 Запуск Jira-сервиса...")

	// Создаем клиента для Jira
	jiraClient := adapter.NewJiraClient(cfg)

	// Создаем сервис Jira
	jiraService := service.NewJiraService(jiraClient)

	// Создаем HTTP-обработчик
	jiraHandler := handler.NewJiraHandler(jiraService)

	// Создаем приложение Fiber
	app := fiber.New()

	// Роутинг (заглушка)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "✅ Jira-service is running"})
	})

	// ✅ Регистрируем маршруты
	app.Get("/issue/:id", jiraHandler.GetIssue)       // Получить одну задачу
	app.Get("/issues", jiraHandler.GetMultipleIssues) // Получить список задач

	// Запускаем сервер
	logger.Info().Msgf("🚀 Сервис запущен на порту %s", cfg.ServerPort)
	err := app.Listen(":" + cfg.ServerPort)
	if err != nil {
		logger.Fatal().Err(err).Msg("❌ Ошибка запуска сервера")
	}
}
