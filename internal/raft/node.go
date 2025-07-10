package raft

import (
	"fmt"
	"sync"
)

type Role int

const (
	Follower Role = iota
	Leader
)

type Node struct {
	ID    int
	Role  Role
	store map[string]string
	mu    sync.RWMutex
	peers []*Node
	inbox chan Message
	// replicas int
}

type Message struct {
	From  int
	Type  string // e.g., "SET"
	Key   string
	Value string
}

func NewNode(id int) *Node {
	return &Node{
		ID:    id,
		Role:  Follower,
		store: make(map[string]string),
		inbox: make(chan Message, 100),
	}
}

func (n *Node) Start() {
	go func() {
		for msg := range n.inbox {
			switch msg.Type {
			case "SET":
				n.mu.Lock()
				n.store[msg.Key] = msg.Value
				n.mu.Unlock()
				fmt.Printf("[Node %d] SET %s=%s from Node %d\n", n.ID, msg.Key, msg.Value, msg.From)
			}
		}
	}()
}

func (n *Node) GetStore() map[string]string {
	n.mu.RLock()
	defer n.mu.RUnlock()
	copy := make(map[string]string)
	for k, v := range n.store {
		copy[k] = v
	}
	return copy
}
