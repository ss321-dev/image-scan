package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ss321-dev/image-scan/internal/code"
	"github.com/ss321-dev/image-scan/internal/conf"
	"github.com/ss321-dev/image-scan/pkg/dockle"
	"github.com/ss321-dev/image-scan/pkg/githubissue"
)

func DockleScan(config conf.Config) (int, error) {

	ctx := context.Background()

	// scan dockle
	dockleConfig := dockle.Config{
		ScanImageName: config.ScanImageName,
		IsLocalImage:  config.IsLocalImage,
	}
	results, err := dockle.Scan(ctx, dockleConfig)
	if err != nil {
		return code.ErrorExitCode, err
	}

	issueSearch := githubissue.IssueSearch{
		Labels: []string{config.IssueApplicationType, config.ImageScanType, config.IssueEnvironment},
	}

	var issueOperations []githubissue.IssueOperation
	exitCode := code.SuccessExitCode
	for _, result := range results.Details {

		if config.IsIssueError(result.Level, result.Code) {
			issueOperations = append(issueOperations, githubissue.IssueOperation{
				Title:  fmt.Sprintf("[%s] %s: %s", result.Level, result.Code, result.Title),
				Body:   strings.Join(result.Alerts, "\n"),
				Labels: issueSearch.Labels,
			})
		}

		if config.IsExitError(result.Level) {
			exitCode = code.ErrorExitCode
		}
	}

	// create issue from the scan results
	gitHubConfig := githubissue.Config{
		AccessToken: config.GitHubAccessToken,
		Owner:       config.GitHubOwner,
		Repository:  config.GitHubRepository,
	}
	if err := githubissue.SaveIssue(ctx, gitHubConfig, issueOperations, issueSearch, true); err != nil {
		if err != nil {
			return code.ErrorExitCode, err
		}
	}
	return exitCode, nil
}
