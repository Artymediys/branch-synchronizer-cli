package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	GitlabURL           string
	GitlabPAT           string
	MattermostURL       string
	MattermostBotToken  string
	MattermostChannelID string
}

func Create(cfg Config) error {
	err := setConfigPath()
	if err != nil {
		return fmt.Errorf("не удаётся получить домашнюю директорию пользователя -> %w", err)
	}

	viper.Set("gitlab_url", cfg.GitlabURL)
	viper.Set("gitlab_pat", cfg.GitlabPAT)
	viper.Set("mm_url", cfg.MattermostURL)
	viper.Set("mm_bot_token", cfg.MattermostBotToken)
	viper.Set("mm_channel_id", cfg.MattermostChannelID)

	if err = viper.WriteConfig(); err != nil {
		return fmt.Errorf("не удаётся записать конфигурационный файл -> %w", err)
	}

	return nil
}

func Read() error {
	err := setConfigPath()
	if err != nil {
		return fmt.Errorf("не удаётся получить домашнюю директорию пользователя -> %w", err)
	}

	if err = viper.ReadInConfig(); err != nil {
		return fmt.Errorf("не удаётся прочитать конфигурационный файл -> %w", err)
	}

	switch {
	case !viper.IsSet("gitlab_url"):
		return fmt.Errorf("не удаётся найти GitLab URL")
	case !viper.IsSet("gitlab_pat"):
		return fmt.Errorf("не удаётся найти GitLab Personal Access Token")
	case !viper.IsSet("mm_url"):
		return fmt.Errorf("не удаётся найти Mattermost URL")
	case !viper.IsSet("mm_bot_token"):
		return fmt.Errorf("не удаётся найти Mattermost Bot Token")
	case !viper.IsSet("mm_channel_id"):
		return fmt.Errorf("не удаётся найти Mattermost Channel ID")
	}

	return nil
}

func setConfigPath() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.SetConfigFile(home + "/.bsync_cli.yaml")

	return nil
}
