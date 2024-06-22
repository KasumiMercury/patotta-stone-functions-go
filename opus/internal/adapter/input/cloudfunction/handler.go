package cloudfunction

import "net/http"

type CloudFunctionHandler struct {
}

func NewCloudFunctionHandler() *CloudFunctionHandler {
	return &CloudFunctionHandler{}
}

func (h *CloudFunctionHandler) Handle(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement

	w.WriteHeader(http.StatusOK)
}
