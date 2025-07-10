package raft

import "fmt"

type Cluster struct {
	Nodes  []*Node
	Leader *Node
}

func NewCluster(size int) *Cluster {
	var nodes []*Node
	for i := 0; i < size; i++ {
		node := NewNode(i)
		nodes = append(nodes, node)
	}

	// Link peers (full mesh)
	for _, n := range nodes {
		n.peers = nodes
		n.Start()
	}

	// Elect static leader (node 0)
	nodes[0].Role = Leader

	return &Cluster{
		Nodes:  nodes,
		Leader: nodes[0],
	}
}

// Send SET to leader, who replicates it to followers
func (c *Cluster) Set(key, val string) {
	leader := c.Leader
	leader.mu.Lock()
	leader.store[key] = val
	leader.mu.Unlock()
	fmt.Printf("[Leader %d] SET %s=%s (self)\n", leader.ID, key, val)

	for _, peer := range leader.peers {
		if peer.ID == leader.ID {
			continue
		}
		peer.inbox <- Message{
			From:  leader.ID,
			Type:  "SET",
			Key:   key,
			Value: val,
		}
	}
}
