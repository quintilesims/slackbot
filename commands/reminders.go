package commands

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/utils"
	"github.com/urfave/cli"
)

const (
	dateFormat = "01/02 2006"
	timeFormat = "03:04PM"
)

// NewRemindersCommand returns a cli.Command that manages !reminders
func NewRemindersCommand(
	store db.Store,
	w io.Writer,
	generateID utils.IDGenerator,
	userParser utils.SlackUserParser,
) cli.Command {
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
						Usage: "the date of the reminder in `mm/dd` format (e.g. `03/15`) (default: <tomorrow>)",
					},
					cli.StringFlag{
						Name:  "time",
						Value: "09:00AM",
						Usage: "the time of the reminder in `HH:MM<AM|PM>` format (e.g. `03:15PM`)",
					},
					cli.IntFlag{
						Name:  "year",
						Usage: "the year of the reminder (e.g. `2015`) (default: <current year>)",
					},
				},
				Action: func(c *cli.Context) error {
					return addReminder(c, store, w, generateID, userParser)
				},
			},
			{
				Name:      "ls",
				Usage:     "list reminders for a user",
				ArgsUsage: "@USER",
				Action: func(c *cli.Context) error {
					return listReminders(c, store, w, userParser)
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

func addReminder(
	c *cli.Context,
	store db.Store,
	w io.Writer,
	generateID utils.IDGenerator,
	userParser utils.SlackUserParser,
) error {
	escapedUser := c.Args().Get(0)
	if escapedUser == "" {
		return fmt.Errorf("@USER is required")
	}

	message := c.Args().Get(1)
	if message == "" {
		return fmt.Errorf("MESSAGE is required")
	}

	if args := c.Args(); len(args) > 2 {
		message = fmt.Sprintf("%s %s", message, strings.Join(args[2:], " "))
	}

	date := c.String("date")
	if date == "" {
		n := time.Now()
		date = fmt.Sprintf("%.2d/%.2d", n.Month(), n.Day()+1)
	}

	year := c.Int("year")
	if year == 0 {
		year = time.Now().Year()
	}

	format := fmt.Sprintf("%s %s", dateFormat, timeFormat)
	input := fmt.Sprintf("%s %.4d %s", date, year, strings.ToUpper(c.String("time")))
	t, err := time.Parse(format, input)
	if err != nil {
		return err
	}

	if t.Before(time.Now()) {
		return fmt.Errorf("Cannot create a reminder in the past!")
	}

	user, err := userParser(escapedUser)
	if err != nil {
		return err
	}

	reminders := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
		return err
	}

	reminderID := generateID()
	reminders[reminderID] = models.Reminder{
		UserID:   user.ID,
		UserName: user.Name,
		Message:  message,
		Time:     time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, time.Local).UTC(),
	}

	log.Printf("[INFO] Added reminder %s", reminders[reminderID])
	if err := store.Write(models.StoreKeyReminders, reminders); err != nil {
		return err
	}

	format = fmt.Sprintf("%s at %s", dateFormat, timeFormat)
	text := fmt.Sprintf("Ok, I've added a new reminder for %s on %s", user.Name, t.Format(format))
	if _, err := w.Write([]byte(text)); err != nil {
		return err
	}

	return nil
}

func listReminders(c *cli.Context, store db.Store, w io.Writer, userParser utils.SlackUserParser) error {
	escapedUser := c.Args().Get(0)
	if escapedUser == "" {
		return fmt.Errorf("@USER is required")
	}

	user, err := userParser(escapedUser)
	if err != nil {
		return err
	}

	reminders := models.Reminders{}
	if err := store.Read(models.StoreKeyReminders, &reminders); err != nil {
		return err
	}

	userReminders := models.Reminders{}
	for reminderID, r := range reminders {
		if r.UserID == user.ID {
			userReminders[reminderID] = r
		}
	}

	if len(userReminders) == 0 {
		text := fmt.Sprintf("%s doesn't have any reminders at the moment", user.Name)
		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}

	text := fmt.Sprintf("%s has the following reminders:\n", user.Name)
	for reminderID, r := range userReminders {
		format := fmt.Sprintf("%s on %s", timeFormat, dateFormat)
		text += fmt.Sprintf("Reminder `%s`: %s at %s\n", reminderID, r.Message, r.Time.Format(format))
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
