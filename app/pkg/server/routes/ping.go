package routes

import (
	"fmt"
	"log"
	"net/http"
)

func (h Handler) Ping(
	w http.ResponseWriter,
	r *http.Request,
) {
	_, err := fmt.Fprintf(w, "pong")
	if err != nil {
		log.Printf("error handling ping, url=%s, err=%s", r.URL, err.Error())
	}
}
