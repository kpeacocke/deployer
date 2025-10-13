package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitHubClient_GetLatestRelease(t *testing.T) {
	// Mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/test/repo/releases/latest" {
			t.Errorf("Expected path '/repos/test/repo/releases/latest', got '%s'", r.URL.Path)
		}

		if r.Header.Get("Accept") != "application/vnd.github.v3+json" {
			t.Errorf("Expected Accept header 'application/vnd.github.v3+json', got '%s'", r.Header.Get("Accept"))
		}

		response := `{
			"tag_name": "v1.2.3",
			"assets": [
				{
					"name": "app.tar.gz",
					"browser_download_url": "https://github.com/test/repo/releases/download/v1.2.3/app.tar.gz"
				},
				{
					"name": "app.zip",
					"browser_download_url": "https://github.com/test/repo/releases/download/v1.2.3/app.zip"
				}
			]
		}`

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	// Create client with mock server
	client := NewGitHubClient("")
	client.client.Transport = &mockTransport{server: server}

	release, err := client.GetLatestRelease(context.Background(), "test/repo")
	if err != nil {
		t.Fatalf("Failed to get latest release: %v", err)
	}

	if release.TagName != "v1.2.3" {
		t.Errorf("Expected tag name 'v1.2.3', got '%s'", release.TagName)
	}

	if len(release.Assets) != 2 {
		t.Errorf("Expected 2 assets, got %d", len(release.Assets))
	}
}

func TestGitHubClient_GetLatestRelease_WithAuth(t *testing.T) {
	// Mock GitHub API server that checks for auth
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "token test-token" {
			t.Errorf("Expected Authorization header 'token test-token', got '%s'", auth)
		}

		response := `{"tag_name": "v1.0.0", "assets": []}`
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	client := NewGitHubClient("test-token")
	client.client.Transport = &mockTransport{server: server}

	_, err := client.GetLatestRelease(context.Background(), "test/repo")
	if err != nil {
		t.Fatalf("Failed to get latest release: %v", err)
	}
}

func TestRelease_FindAssetWithSuffix(t *testing.T) {
	release := &Release{
		TagName: "v1.0.0",
		Assets: []Asset{
			{Name: "app.tar.gz", BrowserDownloadURL: "https://example.com/app.tar.gz"},
			{Name: "app.zip", BrowserDownloadURL: "https://example.com/app.zip"},
			{Name: "checksums.txt", BrowserDownloadURL: "https://example.com/checksums.txt"},
		},
	}

	// Test finding existing asset
	asset, err := release.FindAssetWithSuffix(".tar.gz")
	if err != nil {
		t.Fatalf("Failed to find asset: %v", err)
	}

	if asset.Name != "app.tar.gz" {
		t.Errorf("Expected asset name 'app.tar.gz', got '%s'", asset.Name)
	}

	// Test finding non-existent asset
	_, err = release.FindAssetWithSuffix(".deb")
	if err == nil {
		t.Error("Expected error when finding non-existent asset")
	}
}

// mockTransport is a custom transport for testing HTTP clients
type mockTransport struct {
	server *httptest.Server
}

func (t *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Redirect all requests to our test server
	req.URL.Scheme = "http"
	req.URL.Host = t.server.URL[7:] // Remove "http://" prefix
	return http.DefaultTransport.RoundTrip(req)
}
