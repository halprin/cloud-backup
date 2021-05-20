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

	installAction := cli.NewCommand("install", "Install automatic backup agent").
		WithArg(cli.NewArg("config file", "The configuration file that describes how and what to backup")).
		WithOption(cli.NewOption("month", "Specify a month for when to backup").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("day", "Specify a day in the month for when to backup").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("weekday", "Specify a day of the week for when to backup").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("hour", "Specify an hour for when to backup").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("minute", "Specify a minute for when to backup").WithType(cli.TypeInt)).
		WithAction(func(args []string, options map[string]string) int {
			log.Println(args)
			log.Println(options)

			return 0
		})

	cliApplication := cli.New("Backup files to the cloud").
		WithCommand(backupAction).
		WithCommand(restoreAction).
		WithCommand(installAction)

	os.Exit(cliApplication.Run(os.Args, os.Stdout))
}
