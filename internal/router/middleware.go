package router

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
)

type trackRequest struct{
	Resource string `json:"url"`
}

var urlRegexp = regexp.MustCompile(`^[-a-zA-Z0-9@:%._\\+~#?&\/=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%._\\+~#?&\/\/=]*)$`)

// valuesToCtx adds values from request message to request ctx.
// If failes to read values from body returns an error.
func (r *router) MwValuesToCtx(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		dec := json.NewDecoder(r.Body)
		var message trackRequest
		err := dec.Decode(&message)
		if err != nil{
			http.Error(
				w, 
				fmt.Sprintf("failed to decode message: %v", err),
				http.StatusBadRequest,
			)
						
			return 
		}

		resource := message.Resource

		if !urlRegexp.MatchString(resource){
			http.Error(
				w, 
				fmt.Sprintf("requested resource is not a url: %v", resource),
				http.StatusBadRequest,
			)
			
			return 
		}

		r = r.WithContext(context.WithValue(r.Context(), keyResource, resource))

		next.ServeHTTP(w, r)
	})
}