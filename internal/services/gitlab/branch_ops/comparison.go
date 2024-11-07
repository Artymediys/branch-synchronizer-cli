package branch_ops

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

func CompareBranches(glc *gitlab.Client, projectID, sourceBranch, targetBranch string) (bool, error) {
	baseCompare, _, err := glc.Repositories.Compare(projectID, &gitlab.CompareOptions{
		From: &targetBranch,
		To:   &sourceBranch,
	})
	if err != nil {
		return false, fmt.Errorf("не удалось сравнить ветки в проекте -> %w\n", err)
	}

	if len(baseCompare.Diffs) > 0 {
		return true, nil
	}

	return false, nil
}
