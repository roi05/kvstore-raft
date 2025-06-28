package store

import (
	"bufio"
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

func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Optional: append to WAL (if you want durability for deletes)
	entry := fmt.Sprintf("DELETE %s\n", key)
	_, _ = s.log.WriteString(entry)

	delete(s.data, key)
}

// Replay WAL on startup to rebuild the in-memory store
func (s *Store) replayLog() {
	// Go back to the beginning of the file
	_, err := s.log.Seek(0, 0)
	if err != nil {
		fmt.Println("seek error:", err)
		return
	}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(s.log)
	for scanner.Scan() {
		line := scanner.Text()
		var key, val string

		if _, err := fmt.Sscanf(line, "SET %s %s", &key, &val); err == nil {
			s.data[key] = val
		} else if _, err := fmt.Sscanf(line, "DELETE %s", &key); err == nil {
			delete(s.data, key)
		} else {
			fmt.Println("unknown log line:", line)
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
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
