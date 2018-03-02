package commands

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/utils"
	"github.com/urfave/cli"
)

const (
	DateFormat = "01/02"
	TimeFormat = "03:04PM"
)

func NewRemindersCommand(store db.Store, w io.Writer, newID func() string) cli.Command {
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
						Usage: "the date of the reminder in `mm/dd` format (e.g. `03/15`)",
					},
					cli.StringFlag{
						Name:  "time",
						Value: "09:00AM",
						Usage: "the time of the reminder in `HH:MM<AM|PM>` format (e.g. `03:15PM`)",
					},
				},
				Action: func(c *cli.Context) error {
					return addReminder(c, store, w, newID)
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

func addReminder(c *cli.Context, store db.Store, w io.Writer, newID func() string) error {
	escapedUser := c.Args().Get(0)
	if escapedUser == "" {
		return fmt.Errorf("@USER is required")
	}

	userID, err := utils.ParseSlackUser(escapedUser)
	if err != nil {
		return err
	}

	message := c.Args().Get(1)
	if message == "" {
		return fmt.Errorf("MESSAGE is required")
	}

	date := c.String("date")
	if date == "today" {
		n := time.Now()
		date = fmt.Sprintf("%.2d/%.2d", n.Day(), n.Month())
	}

	format := fmt.Sprintf("%s %s", DateFormat, TimeFormat)
	input := fmt.Sprintf("%s %s", date, strings.ToUpper(c.String("time")))
	t, err := time.Parse(format, input)
	if err != nil {
		return err
	}

	reminders := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
		return err
	}

	reminderID := newID()
	reminders[reminderID] = models.Reminder{
		UserID:  userID,
		Message: message,
		Time:    t,
	}

	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		return err
	}

	format = fmt.Sprintf("%s at %s", DateFormat, TimeFormat)
	text := fmt.Sprintf("Ok, I've added a new reminder for the specified user on %s", t.Format(format))
	if _, err := w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}

func listReminders(c *cli.Context, store db.Store, w io.Writer) error {
	escapedUser := c.Args().Get(0)
	if escapedUser == "" {
		return fmt.Errorf("@USER is required")
	}

	userID, err := utils.ParseSlackUser(escapedUser)
	if err != nil {
		return err
	}

	reminders := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
		return err
	}

	userReminders := models.Reminders{}
	for reminderID, r := range reminders {
		if r.UserID == userID {
			userReminders[reminderID] = r
		}
	}

	if len(userReminders) == 0 {
		text := "That user doesn't have any reminders at the moment"
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}

	text := "That user has the following reminders:\n"
	for reminderID, r := range userReminders {
		dateTime := r.Time.Format(fmt.Sprintf("%s on %s", TimeFormat, DateFormat))
		text += fmt.Sprintf("Reminder `%s`: %s at %s\n", reminderID, r.Message, dateTime)
	}

	if _, err := w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}

func removeReminder(c *cli.Context, store db.Store, w io.Writer) error {
	reminderID := c.Args().Get(0)
	if reminderID == "" {
		return fmt.Errorf("REMINDER_ID is required")
	}

	reminders := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
		return err
	}

	if _, ok := reminders[reminderID]; !ok {
		return fmt.Errorf("Reminder '%s' does not exist", reminderID)
	}

	delete(reminders, reminderID)
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		return err
	}

	text := fmt.Sprintf("Reminder '%s' has been removed", reminderID)
	if _, err := w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}
