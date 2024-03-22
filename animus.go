package patotta_stone_functions_go

import (
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"net/http"
)

func init() {
	functions.HTTP("Animus", animus)
}

func animus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! from Animus!")
}
