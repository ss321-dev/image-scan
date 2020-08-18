package main

import (
	"context"
	"fmt"
	"os"

	"github.com/dockle_cmd/app/dockle"
)

func main() {

	// TODO: Get it from an environment variable
	dockleConfig := dockle.Config{
		ScanImageName: "xxxxxxx",
		IsLocalImage:  true,
	}

	// scan dockle
	result, err := dockle.Scan(context.Background(), dockleConfig)
	if err != nil {
		os.Exit(1)
	}
	fmt.Println(result) // TODO print for test

	// create issue from the scan results

	// slack notification of scan results

	os.Exit(0)
}
