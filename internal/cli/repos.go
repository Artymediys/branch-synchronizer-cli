package cli

import (
	"branch-synchronizer-cli/internal/services/gitlab/repo_ops"

	"github.com/charmbracelet/huh"
	"github.com/xanzy/go-gitlab"
)

func AskForProjects(
	glc *gitlab.Client,
	group *string,
	projectNames *[]string,
) (*huh.Group, error) {

	groups, err := repo_ops.GetGroupNames(glc)
	if err != nil {
		return nil, err
	}

	projects, err := repo_ops.GetGroupProjectMap(glc)
	if err != nil {
		return nil, err
	}

	return huh.NewGroup(
		huh.NewSelect[string]().
			Title("Выберите группу").
			Value(group).
			Height(8).
			Options(huh.NewOptions(groups...)...),
		huh.NewMultiSelect[string]().
			Title("Выберите проекты/репозитории").
			Value(projectNames).
			Height(12).
			OptionsFunc(func() []huh.Option[string] {
				groupProjects := projects[*group]
				return huh.NewOptions(groupProjects...)
			}, &group).
			Filterable(true),
	), nil
}
