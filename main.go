package main

import (
	"log"
	"os"

	"github.com/ss321-dev/image-scan/internal/conf"
	"github.com/ss321-dev/image-scan/internal/service"
)

func main() {

	// convert environment variable to config
	config := &conf.Config{}
	if err := config.Set(); err != nil {
		log.Fatal(err)
	}

	// image scan and operation git hub issue
	if config.ImageScanType == conf.DockleImageScanType {

		exitCode, err := service.DockleScan(*config)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(exitCode)

	} else if config.ImageScanType == conf.EcrImageScanType {

		exitCode, err := service.EcrImageScan(*config)
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(exitCode)
	}
}
