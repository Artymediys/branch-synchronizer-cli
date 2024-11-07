package cli

import (
	"fmt"

	"github.com/charmbracelet/huh"
)

func AskForAcknowledgement(
	confirm *bool,
	projectNames *[]string,
	group, sourceBranch, targetBranch *string,
) *huh.Group {
	var projectNamesMessage string

	for _, projectName := range *projectNames {
		projectNamesMessage += fmt.Sprintln(projectName)
	}

	if len(*projectNames) > 0 {
		ackText := fmt.Sprintf(
			"СОЗДАЁМ МРы ДЛЯ ЭТИХ ПРОЕКТОВ?\nГруппа: %s\nВетки: %s -> %s\nПроекты:\n%s",
			*group, *sourceBranch, *targetBranch, projectNamesMessage,
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
