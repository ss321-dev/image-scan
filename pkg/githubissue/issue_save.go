package githubissue

import (
	"context"
	"errors"
	"fmt"

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

	openListByRepoOptions := issueSearch.convertIssueListByRepoOptions()
	openListByRepoOptions.State = StateOpen
	openIssues, _, err := client.Issues.ListByRepo(ctx, config.Owner, config.Repository, &openListByRepoOptions)
	if err != nil {
		repositoryName := config.Owner + config.Repository
		return fmt.Errorf("failed to search Issue from repository[%s]: %s", repositoryName, err)
	}

	closeListByRepoOptions := issueSearch.convertIssueListByRepoOptions()
	closeListByRepoOptions.State = StateClose
	closedIssues, response, err := client.Issues.ListByRepo(ctx, config.Owner, config.Repository, &closeListByRepoOptions)
	if err != nil {
		repositoryName := config.Owner + config.Repository
		return fmt.Errorf("failed to search Issue from repository[%s]: %s", repositoryName, err)
	}

	// git hub api limit value
	remainingAccess := response.Rate.Remaining
	resetTime := response.Rate.Reset

	// create issue request
	requestIssueOperations := convertIssuesToIssueOperations(openIssues, StateClose)
	closedIssueOperations := convertIssuesToIssueOperations(closedIssues, StateClose)
	for _, issueOperation := range issueOperations {

		// already open issue
		isOpenIssue := false
		for index, requestIssueOperation := range requestIssueOperations {

			if requestIssueOperation.Title == issueOperation.Title {
				isOpenIssue = true
				requestIssueOperations[index].Body = issueOperation.Body
				requestIssueOperations[index].Labels = issueOperation.Labels
				requestIssueOperations[index].state = StateOpen
				break
			}
		}
		if isOpenIssue {
			continue
		}

		// reoccurring closed issue
		for _, closedIssueOperation := range closedIssueOperations {

			if closedIssueOperation.Title == issueOperation.Title {
				issueOperation.number = closedIssueOperation.number
				break
			}
		}
		requestIssueOperations = append(requestIssueOperations, issueOperation)
	}

	// limit check
	if len(requestIssueOperations) > remainingAccess {
		err := errors.New(fmt.Sprintf("GitHub API rate limit exceeded. This limit resets at %s", resetTime.String()))
		return fmt.Errorf("failed to access API: %s", err)
	}

	// save issue
	for _, requestIssueOperation := range requestIssueOperations {

		issueRequest := requestIssueOperation.convertIssueRequest()
		if requestIssueOperation.number == 0 {
			if _, _, err := client.Issues.Create(ctx, config.Owner, config.Repository, &issueRequest); err != nil {
				return fmt.Errorf("failed to create Issue: %s", err)
			}
		} else if *issueRequest.State != StateClose || isAutoClose {
			if _, _, err := client.Issues.Edit(ctx, config.Owner, config.Repository, requestIssueOperation.number, &issueRequest); err != nil {
				return fmt.Errorf("failed to edit Issue number[%d]: %s", requestIssueOperation.number, err)
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
