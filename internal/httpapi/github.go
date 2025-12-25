package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/samokw/foundre-webhook/internal/domain"
	"github.com/samokw/foundre-webhook/internal/github"
	"github.com/samokw/foundre-webhook/internal/preview"
)

func GithubWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	signature := r.Header.Get("X-Hub-Signature-256")
	event := r.Header.Get("X-GitHub-Event")
	deliveryID := r.Header.Get("X-GitHub-Delivery")

	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	body, err := io.ReadAll(r.Body)
	if err != nil {
		var maxbyteerror *http.MaxBytesError
		if errors.As(err, &maxbyteerror) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")
	if secret == "" {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if !github.VerifyGitHubSignature(secret, body, signature) {
		http.Error(w, "invalid signature", http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if event != "pull_request" {
		log.Printf("event=%s delivery=%s bytes=%d (ignored)", event, deliveryID, len(body))
		_, _ = w.Write([]byte("ignored\n"))
		return
	}

	var pre github.PullRequestEvent
	if err := json.Unmarshal(body, &pre); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	action, allowed := github.MapPRAction(pre.Action)
	if !allowed {
		_, _ = w.Write([]byte("ignored"))
		return
	}
	pr := domain.PreviewRequest{
		Repo:   pre.Repository.FullName,
		Number: pre.Number,
		SHA:    pre.PullRequest.Head.SHA,
		Action: action,
	}

	if err := preview.Handle(pr); err != nil {
		log.Printf("preview handle failed: %s", err.Error())
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	log.Printf("number=%d, repo=%s, sha=%s, action=%s", pr.Number, pr.Repo, pr.SHA, pr.Action)

	log.Printf("event=%s delivery=%s bytes=%d", event, deliveryID, len(body))

	_, _ = w.Write([]byte("received\n"))
}
