package main

import (
	"log"
	"net/http"
	"os"

	"github.com/samokw/foundre-webhook/internal/config"
	"github.com/samokw/foundre-webhook/internal/httpapi"
	"github.com/samokw/foundre-webhook/internal/httpapi/middleware"
	"github.com/samokw/foundre-webhook/internal/preview"
)

func main() {
	baseDomain, err := config.PreviewBaseDomain()
	if err != nil {
		log.Fatal(err)
	}
	previewHandler := preview.Handler{BaseDomain: baseDomain}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", httpapi.Health)
	mux.HandleFunc("/github", httpapi.GithubWebhook(previewHandler))
	handler := middleware.Logging(mux)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
