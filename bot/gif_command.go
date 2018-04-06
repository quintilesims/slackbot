package bot

import (
	"fmt"
	"io"
	"math/rand"
	"net/url"
	"strconv"
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
				Usage: "enable explicit content",
			},
			cli.IntFlag{
				Name:  "limit",
				Value: 20,
				Usage: "limit number of results returned (max size 50)",
			},
			cli.BoolFlag{
				Name:  "random",
				Usage: "return a random gif",
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
			query.Set("mediafilter", "minimal")
			query.Set("safesearch", "strict")
			query.Set("limit", strconv.Itoa(c.Int("limit")))

			if c.Bool("explicit") {
				query.Set("safesearch", "off")
			}

			if c.Int("limit") <= 50 {
				query.Set("limit", strconv.Itoa(c.Int("limit")))
			}

			var response TenorSearchResponse
			if err := client.Get("/v1/search", &response, rclient.Query(query)); err != nil {
				return err
			}

			if len(response.Gifs) == 0 {
				return fmt.Errorf("No gifs matching query '%s'", query.Get("q"))
			}

			if c.Bool("random") {
				return write(w, response.Gifs[rand.Intn(len(response.Gifs))].URL)
			}

			return write(w, response.Gifs[0].URL)
		},
	}
}
