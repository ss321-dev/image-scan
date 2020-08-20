package githubissue

import "github.com/google/go-github/v32/github"

const (
	StateOpen  = "open"
	StateClose = "closed"
)

type IssueSearch struct {
	Labels []string
}

func (i *IssueSearch) convertIssueListByRepoOptions() github.IssueListByRepoOptions {
	return github.IssueListByRepoOptions{
		Labels: i.Labels,
	}
}

type IssueOperation struct {
	number int
	state  string
	Title  string
	Body   string
	Labels []string
}

func (i *IssueOperation) convertIssueRequest() github.IssueRequest {

	state := i.state
	if state != StateOpen && state != StateClose {
		state = StateOpen
	}

	return github.IssueRequest{
		Title:  &i.Title,
		Body:   &i.Body,
		Labels: &i.Labels,
		State:  &state,
	}
}

func extractLabelFromIssue(issue github.Issue) []string {
	var labels []string
	for _, label := range issue.Labels {
		labels = append(labels, *label.Name)
	}
	return labels
}
