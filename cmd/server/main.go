package main

import (
	"net/http"

	"github.com/samokw/foundre-webhook/internal/httpapi"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", httpapi.Health)
	http.ListenAndServe(":8080", mux)
}
