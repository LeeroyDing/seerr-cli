package seerr

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("http://localhost:5055", "test-api-key")
	if client.BaseURL != "http://localhost:5055" {
		t.Errorf("expected BaseURL http://localhost:5055, got %s", client.BaseURL)
	}
	if client.APIKey != "test-api-key" {
		t.Errorf("expected APIKey test-api-key, got %s", client.APIKey)
	}
}

func TestGetURL(t *testing.T) {
	client := NewClient("http://localhost:5055/", "key")
	u, err := client.getURL("/api/v1/search")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "http://localhost:5055/api/v1/search"
	if u.String() != expected {
		t.Errorf("expected %s, got %s", expected, u.String())
	}
}

func TestSearch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != "test-key" {
			t.Error("missing X-Api-Key header")
		}
		if r.Header.Get("X-Requested-With") != "XMLHttpRequest" {
			t.Error("missing X-Requested-With header")
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"results": [{"id": 1, "title": "Test Movie"}]}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-key")
	resp, err := client.Search("query")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}

	if len(resp.Results) != 1 {
		t.Errorf("expected 1 result, got %d", len(resp.Results))
	}
	if resp.Results[0].Title != "Test Movie" {
		t.Errorf("expected Test Movie, got %s", resp.Results[0].Title)
	}
}

func TestRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := NewClient(server.URL, "key")
	err := client.Request(123, "movie")
	if err != nil {
		t.Fatalf("Request failed: %v", err)
	}
}

func TestGetMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"displayName": "Leeroy", "requestCount": 5}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "key")
	user, err := client.GetMe()
	if err != nil {
		t.Fatalf("GetMe failed: %v", err)
	}

	if user.DisplayName != "Leeroy" {
		t.Errorf("expected Leeroy, got %s", user.DisplayName)
	}
	if user.RequestCount != 5 {
		t.Errorf("expected 5, got %d", user.RequestCount)
	}
}
