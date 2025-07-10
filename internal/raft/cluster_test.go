package raft

import (
	"testing"
	"time"
)

func TestClusterSetReplication(t *testing.T) {
	cluster := NewCluster(3)

	// Set key via leader
	cluster.Set("name", "roi")

	// Wait briefly for followers to process messages
	time.Sleep(100 * time.Millisecond)

	for _, node := range cluster.Nodes {
		store := node.GetStore()

		val, ok := store["name"]
		if !ok {
			t.Fatalf("Node %d missing key 'name'", node.ID)
		}
		if val != "roi" {
			t.Fatalf("Node %d expected 'roi', got '%s'", node.ID, val)
		}
	}
}
