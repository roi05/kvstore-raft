package main

import (
	"fmt"
	"time"

	"github.com/roi-05/kvstore-raft/internal/raft"
)

func main() {
	cluster := raft.NewCluster(3)

	cluster.Set("x", "100")
	cluster.Set("y", "200")

	time.Sleep(1 * time.Second) // wait for messages to deliver

	for _, node := range cluster.Nodes {
		fmt.Printf("Node %d store: %+v\n", node.ID, node.GetStore())
	}
}
