package store

import (
	"fmt"
	"os"
	"sync"
)

type Store struct {
	mu   sync.RWMutex
	data map[string]string
	log  *os.File
}

// Initialize the store and open the WAL file
func NewStore() *Store {
	logFile, err := os.OpenFile("store.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	s := &Store{
		data: make(map[string]string),
		log:  logFile,
	}

	s.replayLog()
	return s
}

// Write SET operations to the WAL
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entry := fmt.Sprintf("SET %s %s\n", key, value)
	if _, err := s.log.WriteString(entry); err != nil {
		fmt.Println("log write error:", err)
	}

	s.data[key] = value
}

// Replay WAL on startup to rebuild the in-memory store
func (s *Store) replayLog() {
	buf := make([]byte, 4096)
	n, _ := s.log.Read(buf)
	lines := string(buf[:n])
	for _, line := range splitLines(lines) {
		var cmd, key, val string
		fmt.Sscanf(line, "%s %s %s", &cmd, &key, &val)
		if cmd == "SET" {
			s.data[key] = val
		}
	}
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

func splitLines(s string) []string {
	var lines []string
	var curr string
	for _, c := range s {
		if c == '\n' {
			lines = append(lines, curr)
			curr = ""
		} else {
			curr += string(c)
		}
	}
	if curr != "" {
		lines = append(lines, curr)
	}
	return lines
}
