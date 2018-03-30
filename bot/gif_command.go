package bot

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/urfave/cli"
	"github.com/zpatrick/rclient"
)

// TenorAPIEndpoint for the gif api
const TenorAPIEndpoint = "https://api.tenor.com/"

// TenorSearchResponse is the response type for a Tenor API search
type TenorSearchResponse struct {
	Gifs []Gif `json:"results"`
}

// Gif holds information about a Gif from Tenor
type Gif struct {
	URL string `json:"itemurl"`
}

// NewGIFCommand returns a cli.Command that manages !gif
func NewGIFCommand(endpoint, key string, w io.Writer) cli.Command {
	client := rclient.NewRestClient(endpoint)
	return cli.Command{
		Name:      "!gif",
		Usage:     "display a gif",
		ArgsUsage: "args...",
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "explicit",
				Usage: "This will turn off safe search https://tenor.com/gifapi/documentation#safesearch",
			},
		},
		Action: func(c *cli.Context) error {
			args := c.Args()
			if len(args) == 0 {
				return fmt.Errorf("At least one argument is required")
			}

			query := url.Values{}
			query.Set("key", key)
			query.Set("q", strings.Join(args, " "))
			query.Set("limit", "1")
			if !c.Bool("explicit") {
				query.Set("safesearch", "strict")
			}

			var response TenorSearchResponse
			if err := client.Get("/v1/search", &response, rclient.Query(query)); err != nil {
				return err
			}

			if len(response.Gifs) == 0 {
				return fmt.Errorf("No gifs matching query '%s'", query.Get("q"))
			}

			if _, err := w.Write([]byte(response.Gifs[0].URL)); err != nil {
				return err
			}

			return nil
		},
	}
}
