package bot

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strings"

	"github.com/urfave/cli"
	"github.com/zpatrick/rclient"
)

const GiphyAPIEndpoint = "https://api.giphy.com"

// GiphySearchResponse is the response type for a Giphy API search
type GiphySearchResponse struct {
	Gifs []GiphyGif `json:"data"`
}

// GiphyGif holds information about a Gif from Giphy
type GiphyGif struct {
	URL string `json:"bitly_gif_url"`
}

// NewGIFCommand returns a cli.Command that manages !gif
func NewGIFCommand(endpoint, token string, w io.Writer) cli.Command {
	client := rclient.NewRestClient(endpoint)
	return cli.Command{
		Name:      "!gif",
		Usage:     "display a gif",
		ArgsUsage: "args...",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "rating",
				Value: "pg-13",
				Usage: "the MPAA-style rating for the search, e.g. 'pg', 'pg-13', 'r'",
			},
		},
		Action: func(c *cli.Context) error {
			args := c.Args()
			if len(args) == 0 {
				return fmt.Errorf("At least one argument is required")
			}

			query := url.Values{}
			query.Set("api_key", token)
			query.Set("lang", "en")
			query.Set("rating", c.String("rating"))
			query.Set("q", strings.Join(args, " "))

			var response GiphySearchResponse
			if err := client.Get("/v1/gifs/search", &response, rclient.Query(query)); err != nil {
				return err
			}

			if len(response.Gifs) == 0 {
				return fmt.Errorf("No gifs matching query '%s'", query.Get("q"))
			}

			gif := response.Gifs[rand.Intn(len(response.Gifs))]
			if _, err := w.Write([]byte(gif.URL)); err != nil {
				return err
			}

			return nil
		},
	}
}
