package message

type Message struct {
	From int
	Type string
	Data any
}

type OutMsg struct {
	To  int
	Msg Message
}
