package cli

import (
	"fmt"
	"log"

	"branch-synchronizer-cli/internal/services"
	"branch-synchronizer-cli/internal/services/gitlab/branch_ops"
	"branch-synchronizer-cli/internal/services/gitlab/repo_ops"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/xanzy/go-gitlab"
)

var cliTheme = huh.ThemeBase()

func QA(
	glc *gitlab.Client,
	projectIDs, projectNames *[]string,
	group, sourceBranch, targetBranch *string,
	confirm *bool,
) error {
	var (
		groupForm *huh.Group
		groupErr  error
		spinErr   error
		gitlabErr error
	)

	////////////////////////////////////
	///////// ГРУППЫ И ПРОЕКТЫ /////////
	////////////////////////////////////
	log.Println("получаем данные о проектах...")
	spinErr = spinner.New().
		Title("Получаем данные о проектах...").
		Action(func() {
			groupForm, groupErr = AskForProjects(glc, group, projectNames)
		}).Run()
	if spinErr != nil {
		return fmt.Errorf(ErrorSpinner+"%w", spinErr)
	}

	if groupErr != nil {
		return fmt.Errorf("возникла ошибка при формировании интерфейса для групп и проектов -> %w", groupErr)
	}

	err := huh.NewForm(groupForm).WithTheme(cliTheme).Run()
	if err != nil {
		return fmt.Errorf(ErrorForm+"%w", err)
	}

	if len(*projectNames) <= 0 {
		return fmt.Errorf("возникла ошибка при выборе проектов -> должен быть выбран хотя бы 1 проект")
	}

	////////////////////////////////////
	/////////// СПИСОК ВЕТОК ///////////
	////////////////////////////////////
	log.Println("получаем данные о ветках проектов...")
	spinErr = spinner.New().
		Title("Получаем данные о ветках проектов...").
		Action(func() {
			groupForm, groupErr = AskForBranches(
				glc, projectIDs, projectNames, sourceBranch, targetBranch,
			)
		}).Run()
	if spinErr != nil {
		return fmt.Errorf(ErrorSpinner+"%w", spinErr)
	}

	if groupErr != nil {
		return fmt.Errorf("возникла ошибка при формировании интерфейса для веток проекта -> %w", groupErr)
	}

	err = huh.NewForm(groupForm).WithTheme(cliTheme).Run()
	if err != nil {
		return fmt.Errorf(ErrorForm+"%w", err)
	}

	///////////////////////////////////
	///////// СРАВНЕНИЕ ВЕТОК /////////
	///////////////////////////////////
	var syncedProjects []string
	log.Println(fmt.Sprintf(
		"сравниваем ветки \"%s\" – \"%s\"...",
		*sourceBranch, *targetBranch,
	))
	spinErr = spinner.New().
		Title(fmt.Sprintf(
			"Сравниваем ветки \"%s\" – \"%s\"...",
			*sourceBranch, *targetBranch,
		)).
		Action(func() {
			var isDiffs bool
			for i := len(*projectIDs) - 1; i >= 0; i-- {
				isDiffs, gitlabErr = branch_ops.CompareBranches(glc, (*projectIDs)[i], *sourceBranch, *targetBranch)
				if gitlabErr != nil {
					break
				}
				if isDiffs == false {
					syncedProjects = append(syncedProjects, (*projectNames)[i])

					*projectIDs = append((*projectIDs)[:i], (*projectIDs)[i+1:]...)
					*projectNames = append((*projectNames)[:i], (*projectNames)[i+1:]...)
				}
			}
		}).Run()
	if spinErr != nil {
		return fmt.Errorf(ErrorSpinner+"%w", spinErr)
	}
	if gitlabErr != nil {
		return gitlabErr
	}

	/////////////////////////////////////
	/////////// ПОДТВЕРЖДЕНИЕ ///////////
	/////////////////////////////////////
	log.Println("утверждаем выбор пользователя (будем ли создавать МРы)...")
	err = huh.NewForm(
		AskForAcknowledgement(confirm, projectNames, &syncedProjects, group, sourceBranch, targetBranch),
	).WithTheme(cliTheme).Run()
	if err != nil {
		return fmt.Errorf(ErrorForm+"%w", err)
	}

	if *confirm == false {
		return fmt.Errorf("BSync-CLI: выбор не был утверждён")
	}

	return nil
}

func Action(
	glc *gitlab.Client, mmc *services.MattermostClient, mrURL *string,
	projectID, projectName, sourceBranch, targetBranch string,
) (Stage, error) {
	var (
		stageStatus Stage

		spinErr   error
		gitlabErr error
		mmErr     error
		splitErr  error

		shortName string
		fullName  string
	)

	////////////////////////////////////
	/////////// СОЗДАНИЕ МРа ///////////
	////////////////////////////////////
	log.Println(fmt.Sprintf(
		"создаём МР для веток \"%s\" – \"%s\" в проекте \"%s\"...",
		sourceBranch, targetBranch, projectName,
	))
	spinErr = spinner.New().
		Title(fmt.Sprintf(
			"Создаём МР для веток \"%s\" – \"%s\" в проекте \"%s\"...",
			sourceBranch, targetBranch, projectName,
		)).
		Action(func() {
			*mrURL, gitlabErr = branch_ops.CreateMR(glc, projectID, sourceBranch, targetBranch)
		}).Run()
	if spinErr != nil {
		stageStatus.MergeRequest = -1
		return stageStatus, fmt.Errorf(ErrorSpinner+"%w", spinErr)
	}
	if gitlabErr != nil {
		stageStatus.MergeRequest = -1
		return stageStatus, gitlabErr
	}

	log.Printf("МР создан: %s\n", *mrURL)
	stageStatus.MergeRequest = 1

	//////////////////////////////////////
	//////// ОТПРАВКА НОТИФИКАЦИЙ ////////
	//////////////////////////////////////
	log.Println("отправляем нотификацию в Mattermost...")
	spinErr = spinner.New().
		Title("Отправляем нотификацию...").
		Action(func() {
			shortName, fullName, splitErr = repo_ops.GetSplitProjectName(projectName)
			if splitErr != nil {
				return
			}

			message := fmt.Sprintf(
				"> Created new Merge Request!\nService name: **\"%s\"**\nProject name: **\"%s\"**\n"+
					"Branches: `%s -> %s`\nMerge Request: [link](%s)",
				shortName, fullName, sourceBranch, targetBranch, *mrURL,
			)
			mmErr = mmc.Notify(message)
		}).Run()
	if spinErr != nil {
		stageStatus.Notification = -1
		return stageStatus, fmt.Errorf(ErrorSpinner+"%w", spinErr)
	}

	if splitErr != nil {
		stageStatus.Notification = -1
		return stageStatus, splitErr
	}

	if mmErr != nil {
		stageStatus.Notification = -1
		return stageStatus, mmErr
	}

	log.Println("нотификация отправлена!")
	stageStatus.Notification = 1

	return stageStatus, nil
}
