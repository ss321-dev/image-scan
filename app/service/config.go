package service

import (
	"errors"
	"strings"
)

type Config struct {
	ScanImageName        string   `env:"SCAN_IMAGE,required"`
	IsLocalImage         bool     `env:"IS_LOCAL_IMAGE"`
	ErrorLevel           string   `env:"ERROR_LEVEL" envDefault:"fatal"`
	IssueErrorLevel      string   `env:"ISSUE_ERROR_LEVEL" envDefault:"warn"`
	IgnoreErrorCodes     []string `env:"IGNORE_ERROR_CODES" envSeparator:":"`
	GitHubAccessToken    string   `env:"GIT_HUB_ACCESS_TOKEN,required"`
	GitHubOwner          string   `env:"GIT_HUB_Owner,required"`
	GitHubRepository     string   `env:"GIT_HUB_Repository,required"`
	IssueApplicationType string   `env:"ISSUE_APPLICATION_TYPE,required"`
	IssueScanType        string   `env:"ISSUE_SCAN_TYPE,required"`
	IssueEnvironment     string   `env:"ISSUE_ENVIRONMENT,required"`
}

func (c Config) Validate() error {
	var errorText []string

	if _, ok := ErrorLevel[strings.ToLower(c.ErrorLevel)]; !ok {
		errorText = append(errorText, "ERROR_LEVEL should be set to one of the following values [fatal, warn, info, skip, pass]")
	}

	if _, ok := ErrorLevel[strings.ToLower(c.IssueErrorLevel)]; !ok {
		errorText = append(errorText, "ISSUE_ERROR_LEVEL should be set to one of the following values [fatal, warn, info, skip, pass]")
	}

	if len(errorText) > 0 {
		return errors.New(strings.Join(errorText, "\n"))
	}
	return nil
}

func (c Config) IsError(errorLevel string) bool {
	return convertErrorLevelToNumber(c.ErrorLevel) >= convertErrorLevelToNumber(errorLevel)
}

func (c Config) IsIssueError(errorLevel string, code string) bool {
	for _, ignoreErrorCode := range c.IgnoreErrorCodes {
		if ignoreErrorCode == code {
			return false
		}
	}
	return convertErrorLevelToNumber(c.IssueErrorLevel) >= convertErrorLevelToNumber(errorLevel)
}
