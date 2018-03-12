package commands

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/urfave/cli"
	cache "github.com/zpatrick/go-cache"
)

type GoogleClient struct {
	Client *http.Client
}

type Gif struct {
	ContentURL string
}

// NewGifCommand returns a cli.Command that manages !gif
func NewGifCommand(w io.Writer, client *GoogleClient, cache *cache.Cache, url string) cli.Command {
	// Set Cache
	cache.ClearEvery(time.Hour * 24)

	return cli.Command{
		Name:      "!gif",
		Usage:     "display gif for given search query",
		ArgsUsage: "[args...]",
		Action: func(c *cli.Context) error {
			if searchTerm := strings.Join(c.Args(), "+"); searchTerm != "" {
				// Check Cache
				if gif := cache.Get(searchTerm); gif != nil {
					gif, ok := gif.(*Gif)

					if !ok {
						return fmt.Errorf("Cache error")
					}

					if _, err := w.Write([]byte(gif.ContentURL)); err != nil {
						return err
					}
					return nil
				}

				req, err := newGoogleRequest(url, searchTerm)
				if err != nil {
					return err
				}

				gif, err := lookupGif(w, client, req)
				if gif != nil {
					cache.Add(searchTerm, gif)
				}
				if err != nil {
					return err
				}
			}

			// TODO: Return help gif or funny one
			return fmt.Errorf("Searchquery is required")
		},
	}
}

func lookupGif(w io.Writer, client *GoogleClient, req *http.Request) (*Gif, error) {
	resp, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status: %s, Code: %d", resp.Status, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var regexFindContent = regexp.MustCompile(`var u=\'.*?(http.*?)\'`)
	var regexFindURL = regexp.MustCompile(`http.*\w`)

	matches := regexFindContent.FindAllString(string(bodyBytes), -1)
	if len(matches) <= 0 {
		return nil, fmt.Errorf("No gifs found for search")
	}

	Gifs := make([]*Gif, len(matches))
	for i, match := range matches {
		Gifs[i].ContentURL = regexFindURL.FindString(match)
	}

	if _, err := w.Write([]byte(Gifs[0].ContentURL)); err != nil {
		return nil, err
	}
	return Gifs[0], nil
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
