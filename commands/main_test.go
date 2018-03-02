package commands

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/utils"
	"github.com/urfave/cli"
)

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

func runTestApp(cmd cli.Command, format string, tokens ...interface{}) error {
	app := newTestApp(cmd)
	input := fmt.Sprintf(format, tokens...)
	args := append([]string{""}, utils.ParseShell(input)...)
	return app.Run(args)
}
