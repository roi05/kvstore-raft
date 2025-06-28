package apiserver

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/roi-05/kvstore-raft/internal/store"
)

func TestSetAndGet(t *testing.T) {
	kv := store.NewStore()
	h := NewHandler(kv)

	// ----- Test /set
	req := httptest.NewRequest(http.MethodPost, "/set?key=foo&val=bar", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != "OK" {
		t.Fatalf("expected 'OK', got '%s'", string(body))
	}

	// ----- Test /get
	req = httptest.NewRequest(http.MethodGet, "/get?key=foo", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp = w.Result()
	body, _ = io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	if string(body) != "bar" {
		t.Fatalf("expected 'bar', got '%s'", string(body))
	}
}

func TestGetMissingKey(t *testing.T) {
	kv := store.NewStore()
	h := NewHandler(kv)

	req := httptest.NewRequest(http.MethodGet, "/get?key=missing", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", resp.StatusCode)
	}
}

func TestDeleteKey(t *testing.T) {
	kv := store.NewStore()
	h := NewHandler(kv)

	// First, set a key
	req := httptest.NewRequest(http.MethodPost, "/set?key=hello&val=world", nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	// Now, delete it
	req = httptest.NewRequest(http.MethodDelete, "/delete?key=hello", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", resp.StatusCode)
	}
	if string(body) != "OK" {
		t.Fatalf("expected 'OK', got '%s'", string(body))
	}

	// Try to GET the deleted key
	req = httptest.NewRequest(http.MethodGet, "/get?key=hello", nil)
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)

	resp = w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("expected 404 Not Found, got %d", resp.StatusCode)
	}
}
