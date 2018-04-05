package bot

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/urfave/cli"
	"github.com/zpatrick/rclient"
)

// DatamuseAPIEndpoint for the datamuse api
const DatamuseAPIEndpoint = "https://api.datamuse.com"

// DatamuseResponse is the response type for a Tenor API search
type DatamuseResponse []struct {
	Definitions []string `json:"defs"`
	Word        string   `json:"word"`
}

// NewDefineCommand returns a cli.Command that manages !gif
func NewDefineCommand(endpoint string, w io.Writer) cli.Command {
	client := rclient.NewRestClient(endpoint)
	return cli.Command{
		Name:      "!define",
		Usage:     "display a definition for word",
		ArgsUsage: "WORD OR PHRASE",
		Action: func(c *cli.Context) error {
			args := c.Args()
			if len(args) == 0 {
				return fmt.Errorf("At least one argument is required")
			}

			query := url.Values{}
			query.Set("sp", strings.Join(args, " "))
			query.Set("max", "1")
			query.Set("md", "d")

			var response DatamuseResponse
			if err := client.Get("/words", &response, rclient.Query(query)); err != nil {
				return err
			}

			if len(response) == 0 || len(response[0].Definitions) == 0 {
				return fmt.Errorf("No definitions matching query '%s'", query.Get("sp"))
			}

			text := fmt.Sprintf("Here are the defintions for %s: \n", response[0].Word)
			for i := 0; i < len(response[0].Definitions) && i != 4; i++ {
				text += fmt.Sprintf("*%d* %s \n", i+1, response[0].Definitions[i])
			}

			return write(w, text)
		},
	}
}
