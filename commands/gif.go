package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

// NewGifCommand returns a cli.Command that manages !gif
func NewGifCommand(w io.Writer) cli.Command {
	return cli.Command{
		Name:      "!gif",
		Usage:     "display gif for given search query",
		ArgsUsage: "[args...]",
		Action: func(c *cli.Context) error {
			searchQuery := strings.Join(c.Args(), " ")
			if searchQuery != "" {
				return lookupGif(c, w, searchQuery)
			}

			// TODO: Return help gif or funny one
			return fmt.Errorf("Searchquery is required")
		},
	}
}

func lookupGif(
	c *cli.Context,
	w io.Writer,
	searchQuery string,
) error {
	var url = "https://www.googleapis.com/customsearch/v1?q=" + searchQuery + "&searchType=image&fileType=gif&cx=017761564406976645410:cc4umzwhcbc&key=AIzaSyD_9xeeBrIk_amCdiUv7-H_0S-bLn8oz1k"
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	response := new(models.CustomSearchAPIResponse)
	json.NewDecoder(r.Body).Decode(response)

	// Get First Item
	link := response.Items[0].Link
	// TODO: Verify link is gif

	// rr, ee := http.Get(link)

	// http.DetectContentType()

	// TODO: Cache response
	if _, err := w.Write([]byte(link)); err != nil {
		return err
	}
	return nil
}
