// package main

// import (
// 	"log"
// 	"net/http"
// 	"os"

// 	"github.com/samokw/foundre-webhook/internal/httpapi"
// 	"github.com/samokw/foundre-webhook/internal/httpapi/middleware"
// )

//	func main() {
//		mux := http.NewServeMux()
//		mux.HandleFunc("/health", httpapi.Health)
//		mux.HandleFunc("/github", httpapi.GithubWebhook)
//		handler := middleware.Logging(mux)
//		if err := http.ListenAndServe(":8080", handler); err != nil {
//			log.Printf("server stopped: %v", err)
//			os.Exit(1)
//		}
//	}
package main

import (
	"encoding/json"
	"fmt"
)

type PullRequestEvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	PullRequest struct {
		Draft bool `json:"draft"`
		Head  struct {
			SHA string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`
}

func main() {
	payload := []byte(`{
		"action": "opened",
		"number": 42,
		"pull_request": {
			"draft": true,
			"head": { "sha": "abc123" }
		}
	}`)

	var e PullRequestEvent
	if err := json.Unmarshal(payload, &e); err != nil {
		panic(err)
	}

	fmt.Printf("action=%s number=%d draft=%v sha=%s\n",
		e.Action, e.Number, e.PullRequest.Draft, e.PullRequest.Head.SHA)
}
