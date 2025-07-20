package message

type Message struct {
	From        int
	To          int
	Term        int
	Type        string
	VoteGranted bool // only used for MsgVoteResponse
}

type OutMsg struct {
	To  int
	Msg Message
}

const (
	MsgHeartbeat    = "HEARTBEAT"
	MsgVoteRequest  = "VOTE_REQUEST"
	MsgVoteResponse = "VOTE_RESPONSE"
)
