package service

import (
	"errors"

	"github.com/ss321-dev/image-scan/internal/code"
	"github.com/ss321-dev/image-scan/internal/conf"
)

func EcrImageScan(config conf.Config) (int, error) {
	return code.ErrorExitCode, errors.New("ecr image scan is not yet implemented")
}
