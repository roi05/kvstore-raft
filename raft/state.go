package raft

type Role string

const (
	Follower  Role = "Follower"
	Candidate Role = "Canditate"
	Leader    Role = "Leader"
)
