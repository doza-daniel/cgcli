package main

import (
	"fmt"

	"github.com/doza-daniel/cgcli/cgc"
	"github.com/urfave/cli"
)

var filesCmd = cli.Command{Name: "files"}

var filesListCmd = cli.Command{
	Name: "list",
	Action: func(c *cli.Context) error {
		token := c.GlobalString(tokenFlag.Name)
		projectID := c.String(projectFlag.Name)

		client := cgc.New(token)
		files, err := client.Files(projectID)
		if err != nil {
			return err
		}

		for _, file := range files {
			fmt.Println(file.Name, file.ID)
		}
		return nil
	},
}

var filesUpdateCmd = cli.Command{Name: "update"}
var filesStatCmd = cli.Command{Name: "stat"}
var filesDownloadCmd = cli.Command{Name: "download"}

var projectFlag = cli.StringFlag{Name: "project"}
var fileFlag = cli.StringFlag{Name: "file"}
var destFlag = cli.StringFlag{Name: "dest"}

func init() {
	filesListCmd.Flags = []cli.Flag{projectFlag}
	filesStatCmd.Flags = []cli.Flag{fileFlag}
	filesUpdateCmd.Flags = []cli.Flag{fileFlag}
	filesDownloadCmd.Flags = []cli.Flag{fileFlag, destFlag}

	filesCmd.Subcommands = []cli.Command{
		filesListCmd,
		filesUpdateCmd,
		filesStatCmd,
		filesDownloadCmd,
	}
}
