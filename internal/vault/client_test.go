package vault

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestServer(t *testing.T, responseData map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"data": responseData,
		})
	}))
}

func TestNewClient_MissingAddress(t *testing.T) {
	_, err := NewClient("", "token", "")
	if err == nil {
		t.Fatal("expected error for empty address, got nil")
	}
}

func TestNewClient_MissingToken(t *testing.T) {
	_, err := NewClient("http://127.0.0.1:8200", "", "")
	if err == nil {
		t.Fatal("expected error for empty token, got nil")
	}
}

func TestNewClient_Valid(t *testing.T) {
	c, err := NewClient("http://127.0.0.1:8200", "root", "APP_")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestReadSecrets_NoNamespaceFilter(t *testing.T) {
	server := newTestServer(t, map[string]interface{}{
		"data": map[string]interface{}{
			"DB_HOST": "localhost",
			"DB_PORT": "5432",
		},
	})
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 2 {
		t.Errorf("expected 2 secrets, got %d", len(secrets))
	}
}

func TestReadSecrets_WithNamespaceFilter(t *testing.T) {
	server := newTestServer(t, map[string]interface{}{
		"data": map[string]interface{}{
			"APP_KEY": "secret",
			"OTHER_KEY": "ignored",
		},
	})
	defer server.Close()

	client, err := NewClient(server.URL, "test-token", "APP_")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	secrets, err := client.ReadSecrets("secret/data/myapp")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := secrets["APP_KEY"]; !ok {
		t.Error("expected APP_KEY to be present")
	}
	if _, ok := secrets["OTHER_KEY"]; ok {
		t.Error("expected OTHER_KEY to be filtered out")
	}
}
