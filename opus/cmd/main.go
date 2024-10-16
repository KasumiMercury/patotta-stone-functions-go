package main

import (
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/internal/adapters/cloudfunction"
	"github.com/KasumiMercury/patotta-stone-functions-go/opus/logging"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	// Blank-import the function package so the init() runs
)

var handler *cloudfunction.Handler

func init() {
	handler = cloudfunction.NewCloudFunctionHandler()

	logger := logging.NewCustomLogger()
	slog.SetDefault(logger)

	// Register the function to handle HTTP requests
	functions.HTTP("Opus", EntryPoint)
}

func main() {
	// By default, listen on all interfaces. If testing locally, run with
	// LOCAL_ONLY=true to avoid triggering firewall warnings and
	// exposing the server outside your own machine.
	hostname := ""
	if localOnly := os.Getenv("LOCAL_ONLY"); localOnly == "true" {
		hostname = "127.0.0.1"
	}
	if err := funcframework.StartHostPort(hostname, "8080"); err != nil {
		log.Fatalf("funcframework.StartHostPort: %v\n", err)
	}
}

func EntryPoint(w http.ResponseWriter, r *http.Request) {
	handler.Handle(w, r)
}
