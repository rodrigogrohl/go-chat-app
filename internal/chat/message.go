package chat

import "time"

// Message represents a single Message
type Message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
