package handlers

import "net/http"

type EventHandler interface {
	Post(w http.ResponseWriter, r *http.Request)
}
