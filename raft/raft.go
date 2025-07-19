package raft

import (
	"fmt"
	"sync"
	"time"

	"github.com/roi05/kvstore-raft/kv"
	"github.com/roi05/kvstore-raft/message"
)

type Node struct {
	ID          int        // Unique IDs for this node
	Role        Role       // Current role: Follower, Candidate, Leader
	CurrentTerm int        // Latest term this node has seen
	VotedFor    *int       // Candidate ID that received vote in current term (nil = none)
	Log         []LogEntry // Raft log entries

	CommitIndex int // Index of highest log entry known to be committed
	LastApplied int // Index of highest log entry applied to state machine

	NextIndex  map[int]int // For leader: index of the next log entry to send to each peer
	MatchIndex map[int]int // For leader: highest log entry known to be replicated on each peer

	Peers  []int                // IDs of peer nodes
	Inbox  chan message.Message // Incoming messages (e.g., AppendEntries, RequestVote)
	Outbox chan message.OutMsg  // Messages to be sent to other nodes

	StateMachine *kv.KVStore // Pointer to the key-value store

	ElectionTimer *time.Timer  // Election timeout
	Heartbeat     *time.Ticker // Heartbeat ticker (only for leader)

	mu sync.Mutex // Protects access to node state
}

func (n *Node) Run() {
	for msg := range n.Inbox {
		fmt.Println(msg)
	}
}

func (n *Node) Send(to int, msg message.Message) {
	n.Outbox <- message.OutMsg{To: to, Msg: msg}
}

func NewNode(id int, peers []int) *Node {
	return &Node{ID: id, Peers: peers}
}
