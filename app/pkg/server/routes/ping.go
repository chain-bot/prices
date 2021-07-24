package routes

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func (h Handler) Ping(
	w http.ResponseWriter,
	r *http.Request,
) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
			"url": r.URL,
		}).Errorf("handling ping")
	}
}
