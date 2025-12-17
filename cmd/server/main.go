package main

import (
	"log"
	"net/http"
	"os"

	"github.com/samokw/foundre-webhook/internal/httpapi"
	"github.com/samokw/foundre-webhook/internal/httpapi/middleware"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", httpapi.Health)
	mux.HandleFunc("/github", httpapi.GithubNotify)
	handler := middleware.Logging(mux)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
