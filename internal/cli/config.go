package cli

import (
	"branch-synchronizer-cli/internal/config"

	"github.com/charmbracelet/huh"
)

func AskForConfig(cfg *config.Config) *huh.Group {
	return huh.NewGroup(
		huh.NewInput().
			Title("Введите GitLab URL").
			Placeholder("https://gitlab.example.com").
			Prompt("URL:").
			Value(&cfg.GitlabURL),
		huh.NewInput().
			Title("Введите GitLab PAT").
			Placeholder("Personal Access Token").
			Prompt("PAT:").
			Value(&cfg.GitlabPAT),
		huh.NewInput().
			Title("Введите Mattermost URL").
			Placeholder("https://mattermost.example.team").
			Prompt("URL:").
			Value(&cfg.MattermostURL),
		huh.NewInput().
			Title("Введите Mattermost Bot Token").
			Placeholder("Bot Token").
			Prompt("Token:").
			Value(&cfg.MattermostBotToken),
		huh.NewInput().
			Title("Введите Mattermost Channel ID").
			Placeholder("Channel ID").
			Prompt("ID:").
			Value(&cfg.MattermostChannelID),
	)
}
