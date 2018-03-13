package bot

import (
	"fmt"
	"io/ioutil"
	"log"
	"testing"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/utils"
	"github.com/urfave/cli"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func newMemoryStore(t *testing.T) *db.MemoryStore {
	store := db.NewMemoryStore()
	if err := db.Init(store); err != nil {
		t.Fatal(err)
	}

	return store
}

func newTestApp(cmd cli.Command) *cli.App {
	app := cli.NewApp()
	app.Commands = []cli.Command{cmd}
	app.Writer = ioutil.Discard
	app.ErrWriter = ioutil.Discard
	return app
}

func newSlackMessageEvent(format string, tokens ...interface{}) slack.RTMEvent {
	return slack.RTMEvent{
		Data: &slack.MessageEvent{
			Msg: slack.Msg{Text: fmt.Sprintf(format, tokens...)},
		},
	}
}

func runTestApp(cmd cli.Command, format string, tokens ...interface{}) error {
	app := newTestApp(cmd)
	input := fmt.Sprintf(format, tokens...)
	args, err := utils.ParseShell("slackbot " + input)
	if err != nil {
		return err
	}

	return app.Run(args)
}
