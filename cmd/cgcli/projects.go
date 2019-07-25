package main

import (
	"fmt"

	"github.com/doza-daniel/cgcli/cgc"
	"github.com/urfave/cli"
)

var projectsListCmd = cli.Command{
	Name: "list",
	Action: func(c *cli.Context) error {
		client := cgc.New(c.GlobalString(tokenFlag.Name))

		projects, err := client.Projects()
		if err != nil {
			return err
		}

		for _, p := range projects {
			fmt.Println(p.ID)
		}

		return nil
	},
}
var projectsCmd = cli.Command{Name: "projects"}

func init() {
	projectsCmd.Subcommands = []cli.Command{
		projectsListCmd,
	}
}
