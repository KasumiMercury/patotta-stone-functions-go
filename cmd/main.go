package main

import (
	"log"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	_ "github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	// Blank-import the function package so the init() runs
	_ "github.com/KasumiMercury/patotta-stone-functions-go"
)

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
