package repository

import (
	"github.com/majid-cj/go-chat-server/domain/entity"
)

// MemberRepository ...
type ChatRepository interface {
	AddNewChatMessage(*entity.ChatMessage) error
	ReadChatMessage(string, string) error
	GetChatHistory(string) (entity.ChatMessageHistory, error)
	AddChatRoom(*entity.ChatRoom) error
	GetChatList(string) (entity.ChatList, error)
	GetChatCounter(string) (int64, error)
}
