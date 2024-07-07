package cloudfunction

import "net/http"

type Handler struct {
}

func NewCloudFunctionHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement

	// TODO: RSS service

	// TODO: Update Video Info

	w.WriteHeader(http.StatusOK)
}
