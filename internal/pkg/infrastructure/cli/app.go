package cli

import (
	"github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/infrastructure/dependency_injection"
	"github.com/urfave/cli"
)

func GetApp(version string) *cli.App {
	app := cli.NewApp()

	app.Version = version

	app.Name = "Bookmark"
	app.Usage = ""

	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "elasticsearch-url",
			EnvVar: "ELASTICSEARCH_URL",
			Usage:  "Elasticsearch url",
		},
		cli.StringFlag{
			Name:   "server-port",
			EnvVar: "SERVER_PORT",
			Usage:  "HTTP server port",
		},
	}

	app.Commands = []cli.Command{
		getServerCommand(app.Flags),
	}

	return app
}

func getServerCommand(baseFlags []cli.Flag) cli.Command {
	return cli.Command{
		Name:   "server",
		Action: runServer,
		Usage:  "Run the http server which expose SERP API",
	}
}

func runServer(c *cli.Context) error {
	server, err := dependency_injection.ServiceLocatorInstance.HttpServer()
	if err != nil {
		return err
	}

	err = server.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
