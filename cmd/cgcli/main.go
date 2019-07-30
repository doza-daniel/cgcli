package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "cgcli"
	app.Usage = "CLI tool for accessing CGC Public API."
	app.Version = "1.0.0"

	app.Flags = []cli.Flag{tokenFlag}
	app.Commands = []cli.Command{projectsCmd, filesCmd}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
