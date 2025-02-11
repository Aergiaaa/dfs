package p2p

// Message hold any arbitrary data that
// is being sent over each transport 
// between two nodes. 
type Message struct {
	Payload []byte
}
