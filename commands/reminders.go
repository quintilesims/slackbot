package commands

import (
	"fmt"
	"io"

	"github.com/quintilesims/slackbot/common"
	"github.com/quintilesims/slackbot/db"
	"github.com/urfave/cli"
)

const (
	DateFormat = "01/02"
	TimeFormat = "03:04PM"
)

func NewRemindersCommand(store db.Store, w io.Writer) cli.Command {
	return cli.Command{
		Name:  "!reminders",
		Usage: "operations for reminders",
		Subcommands: []cli.Command{
			{
				Name:      "add",
				Usage:     "add a new reminder",
				ArgsUsage: "@USER MESSAGE",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "date",
						Value: "today",
						Usage: "the date of the reminder in mm/dd format (e.g. 03/15)",
					},
					cli.StringFlag{
						Name:  "time",
						Value: "09:00AM",
						Usage: "the time of the reminder in HH:MM<AM|PM> format (e.g. 03:15PM)",
					},
				},
				Action: func(c *cli.Context) error {
					return addReminder(c, store, w)
				},
			},
			{
				Name:      "ls",
				Usage:     "list reminders for a user",
				ArgsUsage: "@USER",
				Action: func(c *cli.Context) error {
					return listReminders(c, store, w)
				},
			},
			{
				Name:      "rm",
				Usage:     "remove a reminder",
				ArgsUsage: "REMINDER_ID",
				Action: func(c *cli.Context) error {
					return removeReminder(c, store, w)
				},
			},
		},
	}
}

func addReminder(c *cli.Context, store db.Store, w io.Writer) error {
	return fmt.Errorf("Add Reminder not implemented")
}

func listReminders(c *cli.Context, store db.Store, w io.Writer) error {
	return fmt.Errorf("List Reminders not implemented")
}

func removeReminder(c *cli.Context, store db.Store, w io.Writer) error {
	reminderID := c.Args().Get(0)
	if reminderID == "" {
		return fmt.Errorf("REMINDER_ID is required")
	}

	reminders := common.Reminders{}
	if err := store.Read(common.StoreKeyReminders, &reminders); err != nil {
		return err
	}

	if _, ok := reminders[reminderID]; !ok {
		return fmt.Errorf("Reminder '%s' does not exist", reminderID)
	}

	delete(reminders, reminderID)
	if err := store.Write(common.StoreKeyReminders, reminders); err != nil {
		return err
	}

	text := fmt.Sprintf("Reminder '%s' has been removed", reminderID)
	if _, err := w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}
