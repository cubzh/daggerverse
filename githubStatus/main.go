package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	statusAPIURL = "https://api.github.com/repos/%s/%s/statuses/%s"
)

type GithubStatus struct{}

type GithubStatusOpts struct {
	AccessToken Secret
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

func (m *GithubStatus) Post(ctx context.Context, opts GithubStatusOpts) error {

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
		return fmt.Errorf("state (%s) should be \"error\", \"failure\", \"pending\" or \"success\"", status.State)
	}

	if status.Context == "" {
		status.Context = "default"
	}

	// Marshal the status into JSON
	statusBytes, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("unable to marshal status: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(statusBytes))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	token, err := opts.AccessToken.Plaintext(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
