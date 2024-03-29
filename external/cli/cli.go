package cli

import (
	"github.com/halprin/cloud-backup/actions"
	"github.com/halprin/cloud-backup/actions/backupset"
	"github.com/halprin/cloud-backup/actions/restore"
	"github.com/teris-io/cli"
	"log"
	"os"
	"strconv"
)

func Cli() {
	backupAction := cli.NewCommand("backup", "Initiate a backup").
		WithArg(cli.NewArg("config file", "The configuration file that describes how and what to backup")).
		WithAction(func(args []string, options map[string]string) int {
			err := backupset.Backup(args[0])
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
			err := restore.Restore(args[0], args[1], args[2], args[3])
			if err != nil {
				log.Println(err.Error())
				return 1
			}

			return 0
		})

	installAction := cli.NewCommand("install", "Install automatic backup agent").
		WithArg(cli.NewArg("config file", "The configuration file that describes how and what to backup")).
		WithOption(cli.NewOption("month", "Specify a month for when to backup (1 - 12)").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("day", "Specify a day in the month for when to backup (1 - 31)").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("weekday", "Specify a day of the week for when to backup (0 - 7)").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("hour", "Specify an hour for when to backup (0 - 23)").WithType(cli.TypeInt)).
		WithOption(cli.NewOption("minute", "Specify a minute for when to backup (0 - 59)").WithType(cli.TypeInt)).
		WithAction(func(args []string, options map[string]string) int {
			var optionalMonth *int
			var optionalDay *int
			var optionalWeekday *int
			var optionalHour *int
			var optionalMinute *int

			monthString, exists := options["month"]
			if exists {
				month, err := strconv.Atoi(monthString)
				if err != nil {
					log.Println(err.Error())
					return 1
				}

				optionalMonth = &month
			}

			dayString, exists := options["day"]
			if exists {
				day, err := strconv.Atoi(dayString)
				if err != nil {
					log.Println(err.Error())
					return 2
				}

				optionalDay = &day
			}

			weekdayString, exists := options["weekday"]
			if exists {
				weekday, err := strconv.Atoi(weekdayString)
				if err != nil {
					log.Println(err.Error())
					return 3
				}

				optionalWeekday = &weekday
			}

			hourString, exists := options["hour"]
			if exists {
				hour, err := strconv.Atoi(hourString)
				if err != nil {
					log.Println(err.Error())
					return 4
				}

				optionalHour = &hour
			}

			minuteString, exists := options["minute"]
			if exists {
				minute, err := strconv.Atoi(minuteString)
				if err != nil {
					log.Println(err.Error())
					return 5
				}

				optionalMinute = &minute
			}

			err := actions.Install(args[0], optionalMonth, optionalDay, optionalWeekday, optionalHour, optionalMinute)
			if err != nil {
				log.Println(err.Error())
				return 6
			}

			return 0
		})

	uninstallAction := cli.NewCommand("uninstall", "Uninstall the automatic backup agent").
		WithAction(func(args []string, options map[string]string) int {
			err := actions.Uninstall()
			if err != nil {
				log.Println(err.Error())
				return 1
			}

			return 0
		})

	listAction := cli.NewCommand("list", "List backup files that can be restored").
		WithArg(cli.NewArg("config file", "The configuration file that describes how and what to backup")).
		WithArg(cli.NewArg("timestamp", "The timestamp in which to view available backups").AsOptional()).
		WithAction(func(args []string, options map[string]string) int {
			var err error

			if len(args) < 2 {
				err = actions.List(args[0], "")
			} else {
				err = actions.List(args[0], args[1])
			}

			if err != nil {
				log.Println(err.Error())
				return 1
			}

			return 0
		})

	cliApplication := cli.New("Backup files to the cloud").
		WithCommand(backupAction).
		WithCommand(restoreAction).
		WithCommand(installAction).
		WithCommand(uninstallAction).
		WithCommand(listAction)

	os.Exit(cliApplication.Run(os.Args, os.Stdout))
}
