package repo_ops

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func GetGroupProjectMap(glc *gitlab.Client) (map[string][]string, error) {
	groupProjectMap := make(map[string][]string, 100)

	groups, err := getGroupsAndSubgroups(glc)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		projects, err := getProjectsInGroup(glc, group.ID)
		if err != nil {
			return nil, err
		}

		projectInfos := make([]struct {
			ShortName string
			FullName  string
		}, 0, len(projects))

		maxShortNameLen := 0
		for _, project := range projects {
			shortName := getProjectPath(project)
			fullName := project.Name

			if len(shortName) > maxShortNameLen {
				maxShortNameLen = len(shortName)
			}

			projectInfos = append(projectInfos, struct {
				ShortName string
				FullName  string
			}{
				ShortName: shortName,
				FullName:  fullName,
			})
		}

		sort.Slice(projectInfos, func(i, j int) bool {
			return strings.ToLower(projectInfos[i].ShortName) < strings.ToLower(projectInfos[j].ShortName)
		})

		projectNames := make([]string, 0, len(projectInfos))
		for _, project := range projectInfos {
			formattedNames := fmt.Sprintf("%-*s| %s", maxShortNameLen+2, project.ShortName, project.FullName)
			projectNames = append(projectNames, formattedNames)
		}

		groupProjectMap[group.Name] = projectNames
	}

	return groupProjectMap, nil
}

func GetProjectID(glc *gitlab.Client, projectName string) (string, error) {
	_, fullProjectName, err := GetSplitProjectName(projectName)
	if err != nil {
		return "", fmt.Errorf("не удалось получить полное имя проекта -> %w", err)
	}

	projects, _, err := glc.Projects.ListProjects(&gitlab.ListProjectsOptions{
		Search: &fullProjectName,
	})
	if err != nil {
		return "", fmt.Errorf("не удалось получить список проектов -> %w", err)
	}

	for _, project := range projects {
		if strings.TrimSpace(project.Name) == fullProjectName {
			return strconv.Itoa(project.ID), nil
		}
	}

	return "", fmt.Errorf("не удалось найти указанный проект -> %s", fullProjectName)
}

func GetGroupNames(gls *gitlab.Client) ([]string, error) {
	groups, err := getGroupsAndSubgroups(gls)
	if err != nil {
		return nil, err
	}

	groupNames := make([]string, 0, 100)
	for _, group := range groups {
		groupNames = append(groupNames, group.Name)
	}

	sort.Slice(groupNames, func(i, j int) bool {
		return strings.ToLower(groupNames[i]) < strings.ToLower(groupNames[j])
	})

	return groupNames, nil
}

func GetSplitProjectName(projectName string) (string, string, error) {
	nameParts := strings.SplitN(projectName, "|", 2)

	if len(nameParts) == 2 {
		shortName := strings.TrimSpace(nameParts[0])
		fullName := strings.TrimSpace(nameParts[1])

		return shortName, fullName, nil
	}

	return "", "", fmt.Errorf("некорректный формат строки с названиями проекта -> %s", projectName)
}

func getGroupsAndSubgroups(glc *gitlab.Client) ([]*gitlab.Group, error) {
	allGroupsMap := make(map[int]*gitlab.Group, 50)
	opt := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}

	for {
		groups, resp, err := glc.Groups.ListGroups(opt)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить список групп -> %w", err)
		}

		for _, group := range groups {
			if _, exists := allGroupsMap[group.ID]; !exists {
				allGroupsMap[group.ID] = group
				subgroups, err := getSubgroups(glc, group.ID, allGroupsMap)
				if err != nil {
					return nil, fmt.Errorf("не удалось получить подгруппы для группы %s -> %w", group.Name, err)
				}
				for _, subgroup := range subgroups {
					allGroupsMap[subgroup.ID] = subgroup
				}
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	allGroups := make([]*gitlab.Group, 0, len(allGroupsMap))
	for _, group := range allGroupsMap {
		allGroups = append(allGroups, group)
	}

	return allGroups, nil
}

func getSubgroups(glc *gitlab.Client, groupID int, allGroupsMap map[int]*gitlab.Group) ([]*gitlab.Group, error) {
	allSubgroups := make([]*gitlab.Group, 0, 5)
	opt := &gitlab.ListSubGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}

	for {
		subgroups, resp, err := glc.Groups.ListSubGroups(groupID, opt)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить подгруппы -> %w", err)
		}

		for _, subgroup := range subgroups {
			if _, exists := allGroupsMap[subgroup.ID]; !exists {
				allSubgroups = append(allSubgroups, subgroup)
				allGroupsMap[subgroup.ID] = subgroup
				subsubgroups, err := getSubgroups(glc, subgroup.ID, allGroupsMap)
				if err != nil {
					return nil, fmt.Errorf("не удалось получить подгруппы для подгруппы %s -> %w", subgroup.Name, err)
				}
				allSubgroups = append(allSubgroups, subsubgroups...)
			}
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allSubgroups, nil
}

func getProjectsInGroup(glc *gitlab.Client, groupID int) ([]*gitlab.Project, error) {
	allProjects := make([]*gitlab.Project, 0, 5)
	opt := &gitlab.ListGroupProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}

	for {
		projects, resp, err := glc.Groups.ListGroupProjects(groupID, opt)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить список репозиториев в группе -> %w", err)
		}

		allProjects = append(allProjects, projects...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allProjects, nil
}

func getProjectPath(project *gitlab.Project) string {
	urlParts := strings.Split(project.WebURL, "/")
	return urlParts[len(urlParts)-1]
}
