package endpoints

import "net/http"

type Endpoint interface {
	Post(w http.ResponseWriter, r *http.Request)
}
