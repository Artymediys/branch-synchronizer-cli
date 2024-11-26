package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func AskForAcknowledgement(
	confirm *bool,
	unsyncedProject, syncedProjects *[]string,
	group, sourceBranch, targetBranch *string,
) *huh.Group {
	var unsyncedProjectsMessage string
	for _, projectName := range *unsyncedProject {
		unsyncedProjectsMessage += fmt.Sprintln(projectName)
	}

	var syncedProjectsMessage string
	for _, projectName := range *syncedProjects {
		syncedProjectsMessage += fmt.Sprintln(projectName)
	}

	if len(*unsyncedProject) > 0 {
		ackText := fmt.Sprintf(
			"БЫЛИ НАЙДЕНЫ НЕСИНХРОНИЗИРОВАННЫЕ ПРОЕКТЫ\nСОЗДАЁМ МРы ДЛЯ ЭТИХ ПРОЕКТОВ?\n"+
				"Группа: %s\nВетки: %s -> %s\n\nСинхронизированные проекты:\n%s\nНЕсинхронизированные проекты:\n%s\n",
			*group, *sourceBranch, *targetBranch, syncedProjectsMessage, unsyncedProjectsMessage,
		)

		return huh.NewGroup(
			huh.NewConfirm().
				Title(ackText).
				Affirmative("Да, Let's Go!").
				Negative("Нет, Отменяем!").
				Value(confirm),
		)
	}

	return huh.NewGroup(
		huh.NewNote().
			Title("ВСЕ ВЕТКИ УЖЕ СИНХРОНИЗИРОВАНЫ!").
			Description("– \"На этом моя работа окончена\"\n– \"Но ты же ничего не сделал\"\n\n").
			Next(true).
			NextLabel("Нажмите \"Enter\", чтобы завершить программу"),
	)
}
