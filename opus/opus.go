package opus

import (
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"net/http"
)

func init() {
	// Register the function to handle HTTP requests
	functions.HTTP("Opus", opus)
}

func opus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}
