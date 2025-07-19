package transport

import (
	"sync"

	"github.com/roi05/kvstore-raft/raft"
)

type Router struct {
	nodes map[int]*raft.Node
	mu    sync.RWMutex
}

func NewRouter() *Router {
	return &Router{
		nodes: make(map[int]*raft.Node),
	}
}

func (r *Router) Register(node *raft.Node) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.nodes[node.ID] = node
}

func (r *Router) Start() {
	for _, node := range r.nodes {
		n := node // capture loop variable safely
		go func(n *raft.Node) {
			for outboxMsg := range n.Outbox {
				r.mu.Lock()
				receiver, ok := r.nodes[outboxMsg.To]
				r.mu.Unlock()

				if ok {
					receiver.Inbox <- outboxMsg.Msg
				}
			}
		}(n)
	}
}
