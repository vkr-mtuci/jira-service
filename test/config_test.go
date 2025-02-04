package test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vkr-mtuci/jira-service/config"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("SERVER_PORT", "8080")
	os.Setenv("JIRA_BASE_URL", "https://jira.example.com")
	os.Setenv("JIRA_API_URL", "/rest/api/latest/")
	os.Setenv("JIRA_API_TOKEN", "dummy-token")

	cfg := config.LoadConfig()

	assert.Equal(t, "8080", cfg.ServerPort)
	assert.Equal(t, "https://jira.example.com", cfg.JiraBaseURL)
}
