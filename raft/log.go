package raft

type LogEntry struct {
	Term    int
	Command Command
}

type Command struct {
	Op    string // "set", "delete"
	Key   string
	Value string
}
