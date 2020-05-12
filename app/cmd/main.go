package main

import (
	"fmt"
	"os"

	"github.com/holive/feedado/app/feedado"
	"github.com/pkg/errors"
)

func main() {
	app, err := feedado.New()
	if err != nil {
		fmt.Println(errors.Wrap(err, "could not run Feedado").Error())
		os.Exit(1)
	}

	_ = app
}
