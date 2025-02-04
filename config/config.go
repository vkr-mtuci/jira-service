package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	ServerPort   string
	JiraBaseURL  string
	JiraAPIUrl   string
	JiraUsername string
	JiraAPIToken string
}

// LoadConfig загружает переменные окружения в структуру Config
func LoadConfig() *Config {
	envPath := ".env" // Путь относительно папки config
	// Загружаем переменные окружения из .env, если файл есть
	if err := godotenv.Load(envPath); err != nil {
		log.Println("⚠ Нет .env файла, используем переменные окружения")
	}

	// Читаем значения из переменных окружения
	config := &Config{
		ServerPort:   os.Getenv("SERVER_PORT"),
		JiraBaseURL:  os.Getenv("JIRA_BASE_URL"),
		JiraAPIUrl:   os.Getenv("JIRA_API_URL"),
		JiraAPIToken: os.Getenv("JIRA_API_TOKEN"),
	}

	// Проверяем, заданы ли критически важные переменные
	if config.JiraBaseURL == "" || config.JiraAPIUrl == "" || config.JiraAPIToken == "" {
		log.Fatal("❌ Ошибка: Не заданы все обязательные переменные окружения для Jira")
	}

	return config
}
