package entity

type PeerType int

const (
	User PeerType = iota
	Channel
	Chat
)

type Peer struct {
	Username   string
	PeerObject any
	Type       PeerType
}
