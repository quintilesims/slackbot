package utils

import (
	"fmt"
	"testing"

	"github.com/mgutz/str"
)

// todo: we may want/need a way to parse args better than splitting by space
//	https://play.golang.org/p/ztqfYiPSlv
// https://github.com/mattn/go-shellwords
/*

OnMessageEvent(e, w) error {
	args := utils.ParseShell(e.Text)
	if args.Get(0) != "!cmdWithFlags" {
		return nil
	}

	f := flag.NewFlagSet("", 0)
	size := flag.NewStringFlag("size", "large", nil)
	f.Parse(args)
	...

	!! ok - so all I really need is utils.ParseShell()
	since then I can just use f.Args(n) and other stuff if i want

*/
func TestParse(t *testing.T) {
	input := `hello world "it is nice outhere"`

	fmt.Println(str.ToArgv(input))

	// first
}
