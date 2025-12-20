package github

type Repo struct {
	FullName string `json:"full_name"`
}

type PullRequestEvent struct {
	Action      string `json:"action"`
	Number      int    `json:"number"`
	Repository  Repo   `json:"repository"`
	PullRequest struct {
		Draft bool `json:"draft"`
		Head  struct {
			SHA string `json:"sha"`
		} `json:"head"`
	} `json:"pull_request"`
}
