package httpapi

import (
	"fmt"
	"net/http"
)

func GithubNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("notify endpoint called (stub)")

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("received\n"))
}
