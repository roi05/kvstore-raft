package store

import (
	"os"
	"testing"
)

func setupTestStore(t *testing.T) *Store {
	// Create a test WAL file
	testFile := "test_store.log"

	// Remove it if it already exists (clean slate)
	_ = os.Remove(testFile)

	// Redirect the logger to the test WAL
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("failed to create test log file: %v", err)
	}

	// Inject it manually
	s := &Store{
		data: make(map[string]string),
		log:  f,
	}
	return s
}

func teardownTestStore(s *Store, file string) {
	s.log.Close()
	os.Remove(file)
}

func TestSetAndGet(t *testing.T) {
	s := setupTestStore(t)
	defer teardownTestStore(s, "test_store.log")

	s.Set("foo", "bar")

	val, ok := s.Get("foo")
	if !ok {
		t.Fatalf("expected key 'foo' to exist")
	}
	if val != "bar" {
		t.Fatalf("expected 'bar', got '%s'", val)
	}
}

func TestReplayLog(t *testing.T) {
	logFile := "test_store.log"
	{
		s := setupTestStore(t)
		defer s.log.Close()
		s.Set("x", "1")
		s.Set("y", "2")
	}

	// Create a new instance which should replay the log
	s2 := &Store{
		data: make(map[string]string),
	}
	f, err := os.OpenFile(logFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	s2.log = f
	s2.replayLog()
	defer teardownTestStore(s2, logFile)

	val, _ := s2.Get("x")
	if val != "1" {
		t.Fatalf("expected x=1, got %s", val)
	}
}

func TestReplayLogWithDelete(t *testing.T) {
	logFile := "test_store.log"
	{
		s := setupTestStore(t)
		defer s.log.Close()

		s.Set("a", "1")
		s.Delete("a")
		s.Set("b", "2")
	}

	s2 := &Store{
		data: make(map[string]string),
	}
	f, err := os.OpenFile(logFile, os.O_RDWR, 0644)
	if err != nil {
		t.Fatalf("open log: %v", err)
	}
	s2.log = f
	s2.replayLog()
	defer teardownTestStore(s2, logFile)

	_, ok := s2.Get("a")
	if ok {
		t.Fatalf("expected a to be deleted")
	}

	val, ok := s2.Get("b")
	if !ok || val != "2" {
		t.Fatalf("expected b=2, got %s", val)
	}
}
