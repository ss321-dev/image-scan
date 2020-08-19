package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/caarlos0/env"
	"github.com/dockle_cmd/conf"
	"github.com/dockle_cmd/dockle"
	"github.com/dockle_cmd/githubissue"
)

func main() {

	ctx := context.Background()

	// convert environment variable to config
	var config conf.Config
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
		log.Fatal(err)
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

		if config.IsExitError(result.Level) {
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
		log.Fatal(err)
	}

	os.Exit(exitCode)
}
