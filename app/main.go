package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/dockle_cmd/app/dockle"
	"github.com/dockle_cmd/app/githubissue"
	"github.com/dockle_cmd/app/service"
)

func main() {

	ctx := context.Background()

	// convert environment variable to config
	var config service.Config
	if err := env.Parse(&config); err != nil {
		log.Fatal(err)
	}
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}

	// scan dockle
	dockleConfig := dockle.Config{
		ScanImageName: config.ScanImageName,
		IsLocalImage:  config.IsLocalImage,
	}
	results, err := dockle.Scan(ctx, dockleConfig)
	if err != nil {
		os.Exit(1)
	}

	issueSearch := githubissue.IssueSearch{
		Labels: []string{config.IssueApplicationType, config.IssueScanType, config.IssueEnvironment},
	}

	var issueOperations []githubissue.IssueOperation
	exitCode := 0
	for _, result := range results.Details {

		if config.IsIssueError(result.Level, result.Code) {
			issueOperations = append(issueOperations, githubissue.IssueOperation{
				Title:  fmt.Sprintf("[%s] %s: %s", result.Level, result.Code, result.Title),
				Body:   strings.Join(result.Alerts, "\n"),
				Labels: issueSearch.Labels,
			})
		}

		if config.IsError(result.Level) {
			exitCode = 1
		}
	}

	// create issue from the scan results
	gitHubConfig := githubissue.Config{
		AccessToken: config.GitHubAccessToken,
		Owner:       config.GitHubOwner,
		Repository:  config.GitHubRepository,
	}
	if err := githubissue.SaveIssue(ctx, gitHubConfig, issueOperations, issueSearch, true); err != nil {
		os.Exit(1)
	}

	os.Exit(exitCode)
}