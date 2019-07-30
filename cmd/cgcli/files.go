package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/doza-daniel/cgcli/cgc"
	"github.com/urfave/cli"
)

var filesCmd = cli.Command{
	Usage: "A set of commands for manipulating files.",
	Name:  "files",
}

var filesListCmd = cli.Command{
	Name:  "list",
	Usage: fmt.Sprintf("List files that belong under a project provided with '%s' flag.", projectFlag.Name),
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

var filesUpdateCmd = cli.Command{
	Name:      "update",
	Usage:     fmt.Sprintf("Update file that's provided with '%s' flag.", fileFlag.Name),
	UsageText: "Takes the arguments in format 'metadata.key=value' or 'key=value' and updates those fields in a file.",
	Action: func(c *cli.Context) error {

		token := c.GlobalString(tokenFlag.Name)
		fileID := c.String(fileFlag.Name)

		client := cgc.New(token)
		err := client.UpdateFile(fileID, c.Args())
		if err != nil {
			return err
		}
		return nil
	},
}

var filesStatCmd = cli.Command{
	Name: "stat",
	Usage: fmt.Sprintf(
		"Prints a JSON string representing information about a file provided with '%s' flag.",
		fileFlag.Name,
	),
	Action: func(c *cli.Context) error {
		token := c.GlobalString(tokenFlag.Name)
		fileID := c.String(fileFlag.Name)

		client := cgc.New(token)
		file, err := client.StatFile(fileID)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(os.Stdout).Encode(file); err != nil {
			return err
		}

		return nil
	},
}
var filesDownloadCmd = cli.Command{
	Name: "download",
	Usage: fmt.Sprintf(
		"Downloads a file to a destination on the system provided with '%s' flag.",
		destFlag.Name,
	),
	Action: func(c *cli.Context) error {
		token := c.GlobalString(tokenFlag.Name)
		dest := c.String(destFlag.Name)
		fileID := c.String(fileFlag.Name)

		client := cgc.New(token)
		return client.DownloadFile(fileID, dest)
	},
}

var projectFlag = cli.StringFlag{
	Usage: "represents the project ID",
	Name:  "project",
}
var fileFlag = cli.StringFlag{
	Usage: "represents the file ID",
	Name:  "file",
}
var destFlag = cli.StringFlag{
	Usage: "a path on a local system",
	Name:  "dest",
}

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
