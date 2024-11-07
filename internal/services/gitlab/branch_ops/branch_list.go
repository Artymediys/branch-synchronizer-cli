package branch_ops

import (
	"fmt"
	"sort"
	"strings"

	"github.com/xanzy/go-gitlab"
)

func GetCommonBranches(glc *gitlab.Client, projectIDs []string) ([]string, error) {
	branchSets := make([]map[string]struct{}, len(projectIDs))

	for i, projectID := range projectIDs {
		branches, err := getBranchNames(glc, projectID)
		if err != nil {
			return nil, err
		}

		branchSet := make(map[string]struct{}, len(branches)*len(projectIDs))
		for _, branch := range branches {
			branchSet[branch] = struct{}{}
		}
		branchSets[i] = branchSet
	}

	commonBranches := make([]string, 0, 5)
	for branch := range branchSets[0] {
		common := true
		for i := 1; i < len(branchSets); i++ {
			if _, exists := branchSets[i][branch]; !exists {
				common = false
				break
			}
		}
		if common {
			commonBranches = append(commonBranches, branch)
		}
	}

	sort.Slice(commonBranches, func(i, j int) bool {
		branch1, branch2 := commonBranches[i], commonBranches[j]
		branch1HasSlash := strings.Contains(branch1, "/")
		branch2HasSlash := strings.Contains(branch2, "/")

		if branch1HasSlash != branch2HasSlash {
			return !branch1HasSlash
		}

		return strings.ToLower(branch1) < strings.ToLower(branch2)
	})

	return commonBranches, nil
}

func getBranchNames(glc *gitlab.Client, projectID string) ([]string, error) {
	branches, err := getBranches(glc, projectID)
	if err != nil {
		return nil, err
	}

	branchNames := make([]string, 0, len(branches))
	for _, branch := range branches {
		branchNames = append(branchNames, branch.Name)
	}

	return branchNames, nil
}

func getBranches(glc *gitlab.Client, projectID string) ([]*gitlab.Branch, error) {
	var allBranches []*gitlab.Branch
	opt := &gitlab.ListBranchesOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}

	for {
		branches, resp, err := glc.Branches.ListBranches(projectID, opt)
		if err != nil {
			return nil, fmt.Errorf("не удалось получить список веток проекта -> %w", err)
		}

		allBranches = append(allBranches, branches...)

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allBranches, nil
}
