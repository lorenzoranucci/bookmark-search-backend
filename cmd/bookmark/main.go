package main

import (
	"fmt"
	"os"

	"github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/infrastructure/cli"
)

var version = "dev"
var app = cli.GetApp(version)

func main() {
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err.Error())
	}
}
