// Copyright Christian Przybulinski
// All Rights Reserved

//Package endpoints
package endpoints

import "net/http"

//Endpoint is the default interface used to configure all the possible handlers that implement the Post method
type Endpoint interface {
	Post(w http.ResponseWriter, r *http.Request)
}
