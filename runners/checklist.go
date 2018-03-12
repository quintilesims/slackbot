package runners

import (
	"log"

	"github.com/nlopes/slack"
	"github.com/quintilesims/slackbot/db"
	"github.com/quintilesims/slackbot/models"
	"github.com/quintilesims/slackbot/utils"
)

// NewChecklistRunner creates a Runner that sends a DM reminder to any users
// who have at least 1 unchecked item in their checklist
func NewChecklistRunner(store db.Store, client *slack.Client) *Runner {
	return &Runner{
		Name: "Checklist",
		run: func() error {
			checklists := models.Checklists{}
			if err := store.Read(models.StoreKeyChecklists, &checklists); err != nil {
				return err
			}

			errs := []error{}
			for userID, checklist := range checklists {
				var hasUncheckedItem bool
				for _, item := range checklist {
					if !item.IsChecked {
						hasUncheckedItem = true
						break
					}
				}

				if hasUncheckedItem {
					text := "Hi! It looks like you have some unfinished items in your checklist.\n"
					text += "Please use `/checklist` to mark or delete items as you complete them"
					params := slack.NewPostMessageParameters()

					log.Printf("[DEBUG] [ChecklistRunner] Sending reminder to %s", userID)
					if _, _, err := client.PostMessage(userID, text, params); err != nil {
						errs = append(errs, err)
					}
				}
			}

			return utils.MultiError(errs)
		},
	}
}
