package raft

import (
	"fmt"
	"sync"
	"time"

	"github.com/roi05/kvstore-raft/kv"
	"github.com/roi05/kvstore-raft/message"
	"github.com/roi05/kvstore-raft/utils"
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
	for {
		switch n.Role {
		case Candidate:
			n.runCandidate()
		case Follower:
			n.runFollower()
		case Leader:
			n.runLeader()
		default:
			return
		}
	}
}

func (n *Node) Send(to int, msg message.Message) {
	n.Outbox <- message.OutMsg{To: to, Msg: msg}
}

func NewNode(id int, peers []int) *Node {
	return &Node{ID: id, Peers: peers}
}

func (n *Node) runFollower() {
	for {
		duration := utils.ElectionTimeoutDuration()
		n.ElectionTimer = time.NewTimer(duration)

		select {
		case <-n.Inbox:
			fmt.Println("Received message")

			duration := utils.ElectionTimeoutDuration()

			if !n.ElectionTimer.Stop() {
				<-n.ElectionTimer.C
			}
			n.ElectionTimer.Reset(duration)

		case <-n.ElectionTimer.C:
			n.Role = Candidate
			fmt.Println("Election timeout")
			return
		}
	}
}

func (n *Node) runCandidate() {
	n.mu.Lock()
	n.CurrentTerm++
	n.VotedFor = &n.ID
	n.mu.Unlock()

	voteCount := 1
	duration := utils.ElectionTimeoutDuration()

	if !n.ElectionTimer.Stop() {
		select {
		case <-n.ElectionTimer.C:
		default:
		}
	}
	n.ElectionTimer.Reset(duration)

	for _, peer := range n.Peers {
		n.Send(peer, message.Message{
			From: n.ID,
			To:   peer,
			Term: n.CurrentTerm,
			Type: message.MsgVoteRequest,
		})
	}

	for {
		select {
		case msg := <-n.Inbox:
			if msg.Term > n.CurrentTerm {
				n.Role = Follower
				n.CurrentTerm = msg.Term
				n.VotedFor = nil
				go n.runFollower()
				return
			}

			if msg.Type == message.MsgVoteResponse && msg.Term == n.CurrentTerm {
				voteCount++
				if voteCount > len(n.Peers)/2 {
					fmt.Println("[Node", n.ID, "] Won election, becoming Leader")
					n.Role = Leader
					go n.runLeader()
					return
				}
			}

		case <-n.ElectionTimer.C:
			fmt.Println("[Node", n.ID, "] Election timeout, restarting election")
			go n.runCandidate()
			return
		}
	}
}

func (n *Node) runLeader() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			for _, peer := range n.Peers {
				n.Send(peer, message.Message{
					From: n.ID,
					To:   peer,
					Term: n.CurrentTerm,
					Type: message.MsgHeartbeat,
				})
				fmt.Printf("[Leader %d] Sent heartbeat to Node %d\n", n.ID, peer)
			}

		case msg := <-n.Inbox:
			if msg.Term >= n.CurrentTerm {
				fmt.Println("Saw higher term, become follower")
				n.Role = Follower
				n.CurrentTerm = msg.Term
				n.VotedFor = nil

				return // you'll transition to follower outside
			}
		}
	}
}
