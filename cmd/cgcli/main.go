package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "cgcli"
	app.Usage = "cli tool for accessing CGC Public API"
	app.Version = "0.0.1"

	app.Flags = []cli.Flag{tokenFlag}
	app.Commands = []cli.Command{projectsCmd, filesCmd}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
