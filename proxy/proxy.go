package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	functions := strings.Split(os.Getenv("FUNCTIONS"), ",")

	for _, function := range functions {
		fun := function
		http.HandleFunc(fmt.Sprintf("/%s", fun), func(w http.ResponseWriter, r *http.Request) {
			parsed, _ := url.Parse(fmt.Sprintf("http://%s:8080/", fun))
			proxy := httputil.NewSingleHostReverseProxy(parsed)
			proxy.ServeHTTP(w, r)
		})
	}

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
