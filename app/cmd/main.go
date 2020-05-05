package main

import (
	"fmt"
	"os"

	"github.com/holive/feed/app/feed"
	"github.com/pkg/errors"
)

func main() {
	app, err := feed.New()
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not run Feed").Error())
		os.Exit(1)
	}

	_ = app
}
