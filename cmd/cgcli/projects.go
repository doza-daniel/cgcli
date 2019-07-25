package main

import "github.com/urfave/cli"

var projectsListCmd = cli.Command{Name: "list"}
var projectsCmd = cli.Command{Name: "projects"}

func init() {
	projectsCmd.Subcommands = []cli.Command{
		projectsListCmd,
	}
}
