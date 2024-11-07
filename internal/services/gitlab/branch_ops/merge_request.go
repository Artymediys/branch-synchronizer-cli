package branch_ops

import (
	"fmt"

	"github.com/xanzy/go-gitlab"
)

func CreateMR(glc *gitlab.Client, projectID, sourceBranch, targetBranch string) (string, error) {
	mrTitle := fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch)
	mrOptions := &gitlab.CreateMergeRequestOptions{
		SourceBranch: &sourceBranch,
		TargetBranch: &targetBranch,
		Title:        &mrTitle,
	}

	mr, _, err := glc.MergeRequests.CreateMergeRequest(projectID, mrOptions)
	if err != nil {
		return "", err
	}

	return mr.WebURL, nil
}
