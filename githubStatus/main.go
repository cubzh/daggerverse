package main

// echo '{githubStatus{post(accessToken: "foo", state: "pending")}}' | dagger query

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	statusAPIURL = "https://api.github.com/repos/%s/%s/statuses/%s"
)

type GithubStatus struct{}

type GithubStatusOpts struct {
	AccessToken string
	Owner       string
	Repo        string
	Sha         string
	State       string // error, failure, pending, success
	TargetURL   string
	Description string
	Context     string // "default" by default
}

type githubStatus struct {
	State       string `json:"state"`
	TargetURL   string `json:"target_url,omitempty"`
	Description string `json:"description,omitempty"`
	Context     string `json:"context,omitempty"`
}

type GithubStatusOutput struct {
	Success bool `json:"success"`
}

// GithubStatusOutput needs at least one function to become a valid type for GraphQL
func (o *GithubStatusOutput) Banane() {}

func (m *GithubStatus) Post(opts GithubStatusOpts) (*GithubStatusOutput, error) {

	url := fmt.Sprintf(statusAPIURL, opts.Owner, opts.Repo, opts.Sha)

	status := githubStatus{
		State:       opts.State,
		TargetURL:   opts.TargetURL,
		Description: opts.Description,
		Context:     opts.Context,
	}

	validStates := map[string]bool{
		"error":   true,
		"failure": true,
		"pending": true,
		"success": true,
	}

	_, isStateValid := validStates[status.State]
	if isStateValid == false {
		return &GithubStatusOutput{Success: false}, fmt.Errorf("state (%s) should be error, failure, pending or success", status.State)
	}

	if status.Context == "" {
		status.Context = "default"
	}

	// Marshal the status into JSON
	statusBytes, err := json.Marshal(status)
	if err != nil {
		return &GithubStatusOutput{Success: false}, fmt.Errorf("unable to marshal status: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(statusBytes))
	if err != nil {
		return &GithubStatusOutput{Success: false}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+opts.AccessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &GithubStatusOutput{Success: false}, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return &GithubStatusOutput{Success: false}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return &GithubStatusOutput{Success: true}, nil
}
