package app

import (
	"fmt"
	"log"
	"strings"
	"time"

	"branch-synchronizer-cli/internal/cli"
	"branch-synchronizer-cli/internal/config"
	"branch-synchronizer-cli/internal/report"
	"branch-synchronizer-cli/internal/services"
	"branch-synchronizer-cli/pkg/utils"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/viper"
	"github.com/xanzy/go-gitlab"
)

func Run() {
	///////////////////////////////////
	//////// LOGGER / REPORTER ////////
	///////////////////////////////////
	logFile, err := utils.NewLogger()
	if err != nil {
		log.Println("не удалось настроить логирование -> ", err)
		return
	}
	defer logFile.Close()

	reportFile, err := report.NewReporter()
	if err != nil {
		log.Println("не удалось настроить репортера -> ", err)
		fmt.Println("не удалось настроить репортера -> ", err)
		return
	}
	defer reportFile.Close()

	///////////////////////////////////
	////////// CONFIGURATION //////////
	///////////////////////////////////
	err = config.Read()
	if err != nil {
		var cfg config.Config

		err = huh.NewForm(cli.AskForConfig(&cfg)).WithTheme(huh.ThemeBase()).Run()
		if err != nil {
			log.Println(cli.ErrorForm, err)
			fmt.Println(cli.ErrorForm, err)
			return
		}

		err = config.Create(cfg)
		if err != nil {
			log.Println("не удалось создать конфиг -> ", err)
			fmt.Println("не удалось создать конфиг -> ", err)
			return
		}
	}

	log.Println("Создаём GitLab клиент...")
	glc, err := gitlab.NewClient(
		viper.GetString("gitlab_pat"), gitlab.WithBaseURL(viper.GetString("gitlab_url")),
	)
	if err != nil {
		log.Println("не удалось создать клиент для взаимодействия с GitLab API -> ", err)
		fmt.Println("не удалось создать клиент для взаимодействия с GitLab API -> ", err)
		return
	}

	log.Println("Создаём Mattermost клиент...")
	mmc, err := services.NewMattermostClient(
		viper.GetString("mm_url"), viper.GetString("mm_bot_token"), viper.GetString("mm_channel_id"),
	)
	if err != nil {
		log.Println("не удалось создать клиент для взаимодействия с Mattermost -> ", err)
		fmt.Println("не удалось создать клиент для взаимодействия с Mattermost -> ", err)
		return
	}

	///////////////////////////////////
	/////////// APPLICATION ///////////
	///////////////////////////////////
	var (
		group        string
		projectIDs   []string
		projectNames []string
		sourceBranch string
		targetBranch string
		confirm      bool
	)

	err = cli.QA(glc, &projectIDs, &projectNames, &group, &sourceBranch, &targetBranch, &confirm)
	if err != nil {
		log.Println(err)
		fmt.Println(err)
		return
	}

	var (
		resultString string
		mrURL        string
	)
	for i := 0; i < len(projectIDs); i++ {
		projectName := strings.Join(strings.Fields(projectNames[i]), " ")

		appStage, err := cli.Action(glc, mmc, &mrURL, projectIDs[i], projectName, sourceBranch, targetBranch)
		if err != nil {
			log.Println(err)
		}

		resultString += fmt.Sprintf(
			"Рапорт от: %s\nПроект: %s\nСоздание МР - %s\nНотификация – %s\n\n",
			time.Now().Format(time.DateTime), projectName,
			appStage.MergeRequest.Status(), appStage.Notification.Status(),
		)
		_, err = reportFile.WriteString(resultString)
		if err != nil {
			log.Println("не удалось записать рапорт в файл ->", err)
		}

		log.Println("проверяем наличие оставшихся проектов...")
		spinErr := spinner.New().
			Title("Проверяем наличие оставшихся проектов...").
			Action(func() {
				time.Sleep(10 * time.Second)
			}).Run()
		if spinErr != nil {
			log.Println(cli.ErrorSpinner+"%w", spinErr)
		}
	}

	fmt.Println(resultString)

	log.Println("BSync-CLI: моя работа на этом закончена!")
}
