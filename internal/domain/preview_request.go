package domain

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
)

type PreviewAction string

const (
	Create PreviewAction = "create"
	Update PreviewAction = "update"
	Delete PreviewAction = "delete"
)

type PreviewRequest struct {
	Repo   string
	Number int
	SHA    string
	Action PreviewAction
}

func (pr PreviewRequest) Key() string {
	return fmt.Sprintf("%s#%d", pr.Repo, pr.Number)
}

// SafeName returns a DNS-1123 label (k8s-friendly) stable identifier.
// Example: "octocat-hello-world-pr-1347"
func (pr PreviewRequest) SafeName() string {
	const maxLen = 63

	repo := sanitizeDNSLabel(pr.Repo)
	if repo == "" {
		repo = "repo"
	}

	suffix := fmt.Sprintf("-pr-%d", pr.Number)
	name := repo + suffix
	if len(name) <= maxLen {
		return name
	}

	// If too long, truncate repo part and add a short hash to reduce collisions.
	// Format: <truncated-repo>-<hash><suffix>
	hash := shortHash8(pr.Repo) // hash original repo string (not sanitized) for stability
	// space available for repo part: maxLen - len("-") - len(hash) - len(suffix)
	avail := maxLen - 1 - len(hash) - len(suffix)
	if avail < 1 {
		// Worst-case fallback
		return hash + suffix
	}

	tr := repo
	if len(tr) > avail {
		tr = tr[:avail]
		tr = strings.Trim(tr, "-")
		if tr == "" {
			tr = "repo"
		}
	}

	return tr + "-" + hash + suffix
}

// sanitizeDNSLabel converts an arbitrary string into a DNS-1123 label-ish string:
// - lowercase
// - only [a-z0-9-]
// - collapses runs of '-'
// - trims leading/trailing '-'
func sanitizeDNSLabel(s string) string {
	s = strings.ToLower(s)
	var b strings.Builder
	b.Grow(len(s))

	prevDash := false
	for _, r := range s {
		isAZ := r >= 'a' && r <= 'z'
		is09 := r >= '0' && r <= '9'

		switch {
		case isAZ || is09:
			b.WriteRune(r)
			prevDash = false
		default:
			if !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}

	out := strings.Trim(b.String(), "-")
	return out
}

func shortHash8(s string) string {
	sum := sha1.Sum([]byte(s))
	// 8 hex chars = 32 bits of hash, enough to avoid most truncation collisions
	return hex.EncodeToString(sum[:])[:8]
}
