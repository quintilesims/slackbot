package controllers

import (
	"log"
	"net/http"

	"github.com/quintilesims/slackbot/slash"
	"github.com/zpatrick/fireball"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err, ok := err.(*slash.SlackMessageError); ok {
		response, err := fireball.NewJSONResponse(200, err.Msg)
		if err != nil {
			log.Printf("[ERROR] Failed to marshal json error: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response.Write(w, r)
		return
	}

	log.Printf("[ERROR] An unhandled error occured: %v", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}
