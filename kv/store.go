package kv

import "sync"

type KVStore struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewKVStore() *KVStore {
	return &KVStore{
		data: make(map[string]string),
	}
}

func (s *KVStore) Set(k, v string) {
	s.mu.Lock()
}
func (s *KVStore) Get(k string) string {
	return ""
}
