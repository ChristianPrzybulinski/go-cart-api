// Copyright Christian Przybulinski
// All Rights Reserved

//Endpoints package
package endpoints

import "net/http"

//Default Endpoint used to configure all the possible handlers that implement the Post method
type Endpoint interface {
	Post(w http.ResponseWriter, r *http.Request)
}
