package entity

import (
	"time"

	"github.com/majid-cj/go-chat-server/util"
)

// ChatMessage ...
type ChatMessage struct {
	ID        string    `bson:"id" json:"id"`
	ChatId    string    `bson:"chat_id" json:"chat_id"`
	Sender    string    `bson:"sender" json:"sender"`
	Receiver  string    `bson:"receiver" json:"receiver"`
	Message   string    `bson:"message" json:"message"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
}

// ChatRoom ...
type ChatRoom struct {
	ID        string    `bson:"id" json:"id"`
	Sender    string    `bson:"sender" json:"sender"`
	Receiver  []string  `bson:"receiver" json:"receiver"`
	Message   string    `bson:"message" json:"message"`
	IsRead    bool      `bson:"is_read" json:"is_read"`
	CreatedAt time.Time `bson:"created_at" json:"created_at,omitempty"`
}

type RetrieveChatRoom struct {
	ID        string          `bson:"id" json:"id"`
	Sender    string          `bson:"sender" json:"sender"`
	Receiver  []MemberProfile `bson:"receiver" json:"receiver"`
	Message   string          `bson:"message" json:"message"`
	IsRead    bool            `bson:"is_read" json:"is_read"`
	CreatedAt time.Time       `bson:"created_at" json:"created_at,omitempty"`
}

// ChatMessageHistory ...
type ChatMessageHistory []ChatMessage

// ChatList ...
type ChatList []RetrieveChatRoom

// PrepareChatMessage ...
func (chat *ChatMessage) PrepareChatMessage() {
	chat.ID = util.ULID()
	chat.CreatedAt = util.GetTimeNow()
}

// PrepareChatRoom ...
func (room *ChatRoom) PrepareChatRoom() {
	room.ID = util.ULID()
	room.CreatedAt = util.GetTimeNow()
}
