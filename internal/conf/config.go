package conf

import (
	"errors"
	"fmt"
	"strings"

	"github.com/caarlos0/env"
	"github.com/ss321-dev/image-scan/pkg/dockle"
)

type Config struct {
	ScanImageName         string   `env:"SCAN_IMAGE,required"`
	ImageScanType         string   `env:"IMAGE_SCAN_TYPE,required"`
	IsLocalImage          bool     `env:"IS_LOCAL_IMAGE"`
	ExitDockleErrorLevel  string   `env:"EXIT_DOCKLE_ERROR_LEVEL" envDefault:"fatal"`
	IssueDockleErrorLevel string   `env:"ISSUE_DOCKLE_ERROR_LEVEL" envDefault:"warn"`
	IgnoreErrorCodes      []string `env:"IGNORE_ERROR_CODES" envSeparator:":"`
	GitHubAccessToken     string   `env:"GIT_HUB_ACCESS_TOKEN,required"`
	GitHubOwner           string   `env:"GIT_HUB_Owner,required"`
	GitHubRepository      string   `env:"GIT_HUB_Repository,required"`
	IssueApplicationType  string   `env:"ISSUE_APPLICATION_TYPE,required"`
	IssueEnvironment      string   `env:"ISSUE_ENVIRONMENT,required"`
}

func (c *Config) Set() error {
	if err := env.Parse(c); err != nil {
		return err
	}
	if err := c.validate(); err != nil {
		return err
	}
	return nil
}

func (c Config) IsExitError(errorLevel string) bool {
	return dockle.ConvertErrorLevelToNumber(c.ExitDockleErrorLevel) >= dockle.ConvertErrorLevelToNumber(errorLevel)
}

func (c Config) IsIssueError(errorLevel string, code string) bool {
	for _, ignoreErrorCode := range c.IgnoreErrorCodes {
		if ignoreErrorCode == code {
			return false
		}
	}
	return dockle.ConvertErrorLevelToNumber(c.IssueDockleErrorLevel) >= dockle.ConvertErrorLevelToNumber(errorLevel)
}

func (c Config) validate() error {

	var errorText []string

	// validate global env
	if c.ImageScanType != EcrImageScanType && c.ImageScanType != DockleImageScanType {
		errorText = append(errorText, fmt.Sprintf("IMAGE_SCAN_TYPE should be set to one of the following values [%s, %s]", EcrImageScanType, DockleImageScanType))
	}

	// validate dockle env
	if _, ok := dockle.ErrorLevel[strings.ToLower(c.ExitDockleErrorLevel)]; !ok {
		errorText = append(errorText, "EXIT_ERROR_LEVEL should be set to one of the following values [fatal, warn, info, skip, pass]")
	}
	if _, ok := dockle.ErrorLevel[strings.ToLower(c.IssueDockleErrorLevel)]; !ok {
		errorText = append(errorText, "ISSUE_ERROR_LEVEL should be set to one of the following values [fatal, warn, info, skip, pass]")
	}

	if len(errorText) > 0 {
		return errors.New(strings.Join(errorText, "\n"))
	}
	return nil
}
