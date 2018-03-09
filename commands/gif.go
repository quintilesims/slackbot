package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
)

// NewGifCommand returns a cli.Command that manages !gif
func NewGifCommand(w io.Writer, client *GoogleClient, url string) cli.Command {
	return cli.Command{
		Name:      "!gif",
		Usage:     "display gif for given search query",
		ArgsUsage: "[args...]",
		Action: func(c *cli.Context) error {
			if searchTerm := strings.Join(c.Args(), "+"); searchTerm != "" {
				req, err := newGoogleRequest(url, searchTerm)
				if err != nil {
					return err
				}

				return lookupGif(w, client, req)
			}

			// TODO: Return help gif or funny one
			return fmt.Errorf("Searchquery is required")
		},
	}
}

func lookupGif(w io.Writer, client *GoogleClient, req *http.Request) error {
	resp, err := client.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Status: %s, Code: %d", resp.Status, resp.StatusCode)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	var regexFindContent = regexp.MustCompile(`var u=\'.*?(http.*?)\'`)
	var regexFindURL = regexp.MustCompile(`http.*\w`)

	matches := regexFindContent.FindAllString(bodyString, -1)
	if len(matches) <= 0 {
		return fmt.Errorf("No gifs found for search")
	}

	Gifs := make([]models.Gif, len(matches))
	for i, match := range matches {
		Gifs[i].ContentURL = regexFindURL.FindString(match)
	}

	// TODO: Cache response
	if _, err := w.Write([]byte(Gifs[0].ContentURL)); err != nil {
		return err
	}
	return nil
}

type GoogleClient struct {
	Client *http.Client
}

func NewGoogleClient() *GoogleClient {
	return &GoogleClient{
		Client: &http.Client{},
	}
}

func newGoogleRequest(url string, searchTerm string) (*http.Request, error) {
	url = strings.TrimSuffix(url, "/")
	url += "/search?tbm=isch&tbs=itp:animated&q=" + searchTerm

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_0 like Mac OS X; en-us) AppleWebKit/532.9 (KHTML, like Gecko) Versio  n/4.0.5 Mobile/8A293 Safari/6531.22.7")

	return req, nil
}
