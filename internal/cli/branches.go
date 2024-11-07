package cli

import (
	"fmt"

	"branch-synchronizer-cli/internal/services/gitlab/branch_ops"
	"branch-synchronizer-cli/internal/services/gitlab/repo_ops"

	"github.com/charmbracelet/huh"
	"github.com/xanzy/go-gitlab"
)

func AskForBranches(
	glc *gitlab.Client,
	projectIDs *[]string, projectNames *[]string,
	sourceBranch *string, targetBranch *string,
) (*huh.Group, error) {

	*projectIDs = make([]string, len(*projectNames))

	for i, projectName := range *projectNames {
		projectID, err := repo_ops.GetProjectID(glc, projectName)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить ID проекта -> %w", err)
		}
		(*projectIDs)[i] = projectID
	}

	branches, err := branch_ops.GetCommonBranches(glc, *projectIDs)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить общие ветки проектов -> %w", err)
	}

	return huh.NewGroup(
		huh.NewSelect[string]().
			Title("Выберите исходную ветку").
			Value(sourceBranch).
			Height(8).
			Options(huh.NewOptions(branches...)...),
		huh.NewSelect[string]().
			Title("Выберите целевую ветку").
			Value(targetBranch).
			Height(8).
			Options(huh.NewOptions(branches...)...).
			Validate(func(selection string) error {
				if selection == *sourceBranch {
					return fmt.Errorf("Выберите разные ветки!")
				}

				return nil
			}),
	), nil
}
