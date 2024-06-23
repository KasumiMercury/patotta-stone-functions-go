package cloudfunction

import "net/http"

type CloudFunctionHandler struct {
}

func NewCloudFunctionHandler() *CloudFunctionHandler {
	return &CloudFunctionHandler{}
}

func (h *CloudFunctionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement

	// TODO: RSS service

	// TODO: Update Video Info

	w.WriteHeader(http.StatusOK)
}
