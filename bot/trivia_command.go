package bot

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/urfave/cli"
	"github.com/zpatrick/rclient"
)

// endpoint for the trivia api
const TriviaAPIEndpoint = "https://opentdb.com/api.php"

// TriviaResponse is the response type for a trivia api search
type TriviaResponse struct {
	Questions []TriviaQuestion `json:"results"`
}

// TriviaQuestion holds information about a single trivia question
type TriviaQuestion struct {
	Category         string   `json:"category"`
	Question         string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

// NewTriviaCommand returns a cli.Command that manages !trivia
func NewTriviaCommand(store db.Store, endpoint string, w io.Writer) cli.Command {
	client := rclient.NewRestClient(endpoint)
	return cli.Command{
		Name:  "!trivia",
		Usage: "commands related to trivia",
		Subcommands: []cli.Command{
			{
				Name:      "answer",
				Usage:     "answer the current trivia question",
				ArgsUsage: "ANSWER",
				Action:    newTriviaAnswerAction(store, w),
			},
			{
				Name:  "new",
				Usage: "start a new trivia question",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  "difficulty",
						Value: "medium",
						Usage: "the dificulty of the question; can be 'easy', 'medium', or 'hard'",
					},
				},
				Action: newTriviaStartAction(client, store, w),
			},
			{
				Name:      "show",
				Usage:     "show the current trivia question",
				ArgsUsage: " ",
				Action:    newTriviaShowAction(store, w),
			},
		},
	}
}

func newTriviaAnswerAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		answer := strings.Join(c.Args(), " ")
		if answer == "" {
			return fmt.Errorf("ANSWER is required")
		}

		question := models.TriviaQuestion{}
		if err := store.Read(db.TriviaKey, &question); err != nil {
			if _, ok := err.(db.MissingEntryError); ok {
				return fmt.Errorf("There isn't an active trivia question at the moment")
			}

			return err
		}

		text := fmt.Sprintf("Sorry, *%s* is not the correct answer", answer)
		if strings.ToLower(answer) == strings.ToLower(question.CorrectAnswer) {
			text = fmt.Sprintf("*%s* is the correct answer!", answer)
		}

		if _, err := w.Write([]byte(text)); err != nil {
			return err
		}

		return nil
	}
}

func newTriviaStartAction(client *rclient.RestClient, store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		query := url.Values{}
		query.Set("amount", "1")
		query.Set("difficulty", c.String("difficulty"))

		var response TriviaResponse
		if err := client.Get("", &response, rclient.Query(query)); err != nil {
			return err
		}

		if len(response.Questions) == 0 {
			return fmt.Errorf("No trivia questions returned by the api!")
		}

		question := models.TriviaQuestion{
			Question:         response.Questions[0].Question,
			CorrectAnswer:    response.Questions[0].CorrectAnswer,
			IncorrectAnswers: response.Questions[0].IncorrectAnswers,
		}

		if err := store.Write(db.TriviaKey, question); err != nil {
			return err
		}

		if _, err := w.Write([]byte(question.String())); err != nil {
			return err
		}

		return nil
	}
}

func newTriviaShowAction(store db.Store, w io.Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		question := models.TriviaQuestion{}
		if err := store.Read(db.TriviaKey, &question); err != nil {
			if _, ok := err.(db.MissingEntryError); ok {
				return fmt.Errorf("There isn't an active trivia question at the moment")
			}

			return err
		}

		if _, err := w.Write([]byte(question.String())); err != nil {
			return err
		}

		return nil
	}
}
