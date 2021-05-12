package cli

import (
	"github.com/halprin/cloud-backup-go/actions"
	"github.com/teris-io/cli"
	"log"
	"os"
)

func Cli() {
	backupAction := cli.NewCommand("backup", "Initiate a backup").
		WithArg(cli.NewArg("config file", "The configuration file that describes how and what to backup")).
		WithAction(func(args []string, options map[string]string) int {
			err := actions.Backup(args[0])
			if err != nil {
				log.Println(err.Error())
				return 1
			}

			return 0
	})

	restoreAction := cli.NewCommand("restore", "Restore a file").
		WithArg(cli.NewArg("config file", "The configuration file that describes how and what to backup")).
		WithArg(cli.NewArg("timestamp", "The timestamp that the backup was taken at")).
		WithArg(cli.NewArg("back up file", "The file to restore")).
		WithArg(cli.NewArg("restore path", "The location that the file is restored to")).
		WithAction(func(args []string, options map[string]string) int {
			err := actions.Restore(args[0], args[1], args[2], args[3])
			if err != nil {
				log.Println(err.Error())
				return 1
			}

			return 0
		})

	cliApplication := cli.New("Backup files to the cloud").
		WithCommand(backupAction).
		WithCommand(restoreAction)

	os.Exit(cliApplication.Run(os.Args, os.Stdout))
}