package cloudfunction

import "net/http"

type HandlerRepository interface {
	Handle(w http.ResponseWriter, r *http.Request)
}
