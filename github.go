package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GitHubClient handles GitHub API interactions
type GitHubClient struct {
	token         string
	client        *http.Client
	lastCheck     time.Time
	cachedRelease *Release
}

// Release represents a GitHub release
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// Asset represents a GitHub release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// NewGitHubClient creates a new GitHub client
func NewGitHubClient(token string) *GitHubClient {
	return &GitHubClient{
		token: token,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetLatestRelease gets the latest release for a repository
func (c *GitHubClient) GetLatestRelease(ctx context.Context, repo string) (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("rate limited by GitHub API")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned %d: %s", resp.StatusCode, string(body))
	}

	var release Release
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &release, nil
}

// FindAssetWithSuffix finds an asset with the specified suffix
func (r *Release) FindAssetWithSuffix(suffix string) (*Asset, error) {
	for _, asset := range r.Assets {
		if strings.HasSuffix(asset.Name, suffix) {
			return &asset, nil
		}
	}
	return nil, fmt.Errorf("no asset found with suffix %s", suffix)
}

// DownloadAsset downloads an asset to the specified path
func (c *GitHubClient) DownloadAsset(ctx context.Context, asset *Asset, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, "GET", asset.BrowserDownloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	if c.token != "" {
		req.Header.Set("Authorization", "token "+c.token)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// TODO: Implement actual file download and extraction
	// This would involve creating the destination directory,
	// downloading the file, and extracting if it's a tar.gz

	return fmt.Errorf("download implementation pending")
}
