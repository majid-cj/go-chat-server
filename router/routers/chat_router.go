package routers

import (
	"encoding/json"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/majid-cj/go-chat-server/config"
	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/infrastructure/auth"
	"github.com/majid-cj/go-chat-server/util"
	"github.com/olahol/melody"
)

// ChatRouter ...
type ChatRouter struct {
	Config *config.AppConfig
}

// NewChatRouter ...
func NewChatRouter(config *config.AppConfig) *ChatRouter {
	return &ChatRouter{
		Config: config,
	}
}

// SaveChatMessage ...
func (router *ChatRouter) SaveChatMessage(message *entity.ChatMessage, sender string, receiver string, isRead bool) {
	defer router.Config.Wg.Done()
	var room entity.ChatRoom
	room.ID = util.ULID()
	room.Sender = sender
	room.Receiver = []string{receiver}
	room.CreatedAt = util.GetTimeNow()
	room.Message = message.Message
	room.IsRead = isRead
	router.Config.Persistence.Chat.AddNewChatMessage(message)
	router.Config.Persistence.Chat.AddChatRoom(&room)
}

func (router *ChatRouter) ReadChatMessage(sender, receiver string) {
	defer router.Config.Wg.Done()
	router.Config.Persistence.Chat.ReadChatMessage(sender, receiver)
}

// HandleRequest ...
func (router *ChatRouter) HandleRequest(c iris.Context) {
	if auth.URLTokenValid(c.Request()) {
		router.Config.Melody.HandleRequest(c.ResponseWriter(), c.Request())
	}
}

// HandleConnect ...
func (router *ChatRouter) HandleConnect(s *melody.Session) {
	URL := s.Request.URL.Path
	chatId := util.GetChatId(URL, false)
	sender := util.GetURLIds(URL)[0]
	receiver := util.GetURLIds(URL)[1]

	router.Config.Wg.Add(1)
	go router.ReadChatMessage(sender, receiver)
	router.Config.Wg.Wait()

	router.Config.Set(chatId, s)
	history, _ := router.Config.Persistence.Chat.GetChatHistory(chatId)
	sent, _ := json.Marshal(history)
	s.Write(sent)
}

// HandleDisconnect ...
func (router *ChatRouter) HandleDisconnect(s *melody.Session) {
	chatId := util.GetChatId(s.Request.URL.Path, false)
	router.Config.CloseSession(chatId)
}

// HandleClose ...
func (router *ChatRouter) HandleClose(s *melody.Session, code int, reason string) error {
	if code == 69 {
		chatId := util.GetChatId(s.Request.URL.Path, false)
		router.Config.CloseSession(chatId)
		return nil
	}
	return util.GetError("general_error")
}

// HandleMessage ...
func (router *ChatRouter) HandleMessage(s *melody.Session, msg []byte) {
	var message *entity.ChatMessage
	err := json.Unmarshal(msg, &message)
	if err != nil {
		return
	}
	URL := s.Request.URL.Path
	senderChat := util.GetChatId(URL, false)
	receiverChat := util.GetChatId(URL, true)
	message.PrepareChatMessage()
	message.ChatId = senderChat

	router.Config.Wg.Add(1)
	go router.SaveChatMessage(message, message.Sender, message.Receiver, true)
	router.Config.Wg.Wait()

	if sender := router.Config.Get(senderChat); sender != nil {
		sent, _ := json.Marshal(entity.ChatMessageHistory{*message})
		sender.Write(sent)
	}

	message.PrepareChatMessage()
	message.ChatId = receiverChat

	router.Config.Wg.Add(1)
	go router.SaveChatMessage(message, message.Receiver, message.Sender, false)
	router.Config.Wg.Wait()

	if receiver := router.Config.Get(receiverChat); receiver != nil {

		router.Config.Wg.Add(1)
		go router.ReadChatMessage(message.Receiver, message.Sender)
		router.Config.Wg.Wait()

		sent, _ := json.Marshal(entity.ChatMessageHistory{*message})
		receiver.Write(sent)
	}
}

// GetChatList ...
func (router *ChatRouter) GetChatList(c iris.Context) {
	c.ContentType("text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	tick := time.NewTicker(time.Millisecond * 500)
	profile := auth.ExtractTokenClaims(c.Request(), "profile_id")

	for {
		select {
		case <-tick.C:
			rooms, _ := router.Config.Persistence.Chat.GetChatList(profile)
			value, _ := json.Marshal(rooms)
			c.Writef("data: %s\n\n", value)
			c.ResponseWriter().Flush()

		case <-c.Request().Context().Done():
			tick.Stop()
			return
		}
	}
}

// GetChatCounter ...
func (router *ChatRouter) GetChatCounter(c iris.Context) {
	c.ContentType("text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")

	tick := time.NewTicker(time.Millisecond * 500)
	profile := auth.ExtractTokenClaims(c.Request(), "profile_id")

	for {
		select {
		case <-tick.C:
			counter, _ := router.Config.Persistence.Chat.GetChatCounter(profile)
			value, _ := json.Marshal(counter)
			c.Writef("data: %s\n\n", value)
			c.ResponseWriter().Flush()

		case <-c.Request().Context().Done():
			tick.Stop()
			return
		}
	}
}
