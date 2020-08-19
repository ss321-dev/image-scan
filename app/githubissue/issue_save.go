package githubissue

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

func SaveIssue(ctx context.Context, config Config, issueOperations []IssueOperation, issueSearch IssueSearch, isAutoClose bool) error {

	// create auth client
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.AccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	// search issue only open
	issueListByRepoOptions := issueSearch.convertIssueListByRepoOptions()
	issueListByRepoOptions.State = StateOpen
	issues, response, err := client.Issues.ListByRepo(ctx, config.Owner, config.Repository, &issueListByRepoOptions)
	if err != nil {
		repositoryName := config.Owner + config.Repository
		log.Println(fmt.Errorf("failed to search Issue from repository[%s]: %s", repositoryName, err))
		return err
	}

	// git hub api limit value
	remainingAccess := response.Rate.Remaining
	resetTime := response.Rate.Reset

	// create issue request
	requestIssueOperations := convertIssuesToIssueOperations(issues, StateClose)
	for _, issueOperation := range issueOperations {

		isAlreadyIssue := false
		for index, requestIssueOperation := range requestIssueOperations {

			if requestIssueOperation.Title == issueOperation.Title {
				isAlreadyIssue = true
				requestIssueOperations[index].Body = issueOperation.Body
				requestIssueOperations[index].Labels = issueOperation.Labels
				requestIssueOperations[index].state = StateOpen
				break
			}
		}

		if isAlreadyIssue {
			continue
		}
		requestIssueOperations = append(requestIssueOperations, issueOperation)
	}

	// limit check
	if len(requestIssueOperations) > remainingAccess {
		err := errors.New(fmt.Sprintf("GitHub API rate limit exceeded. This limit resets at %s", resetTime.String()))
		log.Println(fmt.Errorf("failed to access API: %s", err))
		return err
	}

	// save issue
	for _, requestIssueOperation := range requestIssueOperations {

		issueRequest := requestIssueOperation.convertIssueRequest()
		if requestIssueOperation.number == 0 {
			if _, _, err := client.Issues.Create(ctx, config.Owner, config.Repository, &issueRequest); err != nil {
				log.Println(fmt.Errorf("failed to create Issue: %s", err))
				return err
			}
		} else if *issueRequest.State != StateClose || isAutoClose {
			if _, _, err := client.Issues.Edit(ctx, config.Owner, config.Repository, requestIssueOperation.number, &issueRequest); err != nil {
				log.Println(fmt.Errorf("failed to edit Issue number[%d]: %s", requestIssueOperation.number, err))
				return err
			}
		}
	}
	return nil
}

func convertIssuesToIssueOperations(issues []*github.Issue, state string) []IssueOperation {

	if state != StateOpen && state != StateClose {
		state = StateOpen
	}

	var issueOperations []IssueOperation
	for _, issue := range issues {
		issueOperations = append(issueOperations, IssueOperation{
			number: *issue.Number,
			Title:  *issue.Title,
			state:  state,
			Labels: extractLabelFromIssue(*issue),
		})
	}
	return issueOperations
}
