package persistence

import (
	"context"

	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/domain/repository"
	"github.com/majid-cj/go-chat-server/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ChatRepository ...
type ChatRepository struct {
	Ctx context.Context
	DB  *mongo.Database
}

// NewChatRepository ...
func NewChatRepository(db *mongo.Database) *ChatRepository {
	return &ChatRepository{
		Ctx: context.Background(),
		DB:  db,
	}
}

var _ repository.ChatRepository = &ChatRepository{}

// AddNewChatMessage ...
func (repo *ChatRepository) AddNewChatMessage(message *entity.ChatMessage) error {
	_, err := repo.DB.Collection(CHAT).InsertOne(repo.Ctx, &message)
	if err != nil {
		return err
	}
	return nil
}

// ReadChatMessage ...
func (repo *ChatRepository) ReadChatMessage(sender, receiver string) error {
	filter := bson.M{"sender": sender, "receiver": bson.M{"$in": []string{receiver}}, "is_read": false}
	update := bson.M{"$set": bson.M{
		"is_read": true,
	}}
	_, err := repo.DB.Collection(CHAT_ROOM).UpdateMany(repo.Ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

// GetChatHistory ...
func (repo *ChatRepository) GetChatHistory(key string) (entity.ChatMessageHistory, error) {
	var messages entity.ChatMessageHistory
	cursor, err := repo.DB.Collection(CHAT).Find(repo.Ctx, bson.M{"chat_id": key})
	if err != nil {
		return nil, err
	}
	err = cursor.All(repo.Ctx, &messages)
	if err != nil {
		return nil, err
	}
	return messages, nil
}

// AddChatRoom ...
func (repo *ChatRepository) AddChatRoom(room *entity.ChatRoom) error {
	filter := bson.M{"sender": room.Sender, "receiver": bson.M{"$in": room.Receiver}}
	update := bson.M{"$set": bson.M{
		"id":         room.ID,
		"sender":     room.Sender,
		"receiver":   room.Receiver,
		"message":    room.Message,
		"is_read":    room.IsRead,
		"created_at": room.CreatedAt,
	}}
	upsert := true
	_, err := repo.DB.Collection(CHAT_ROOM).UpdateOne(repo.Ctx, filter, update, &options.UpdateOptions{
		Upsert: &upsert,
	})
	if err != nil {
		return util.GetError("general_error")
	}
	return nil
}

// GetChatList ...
func (repo *ChatRepository) GetChatList(sender string) (entity.ChatList, error) {
	var chatList entity.ChatList
	match := bson.D{{"$match", bson.M{
		"sender": sender,
	}}}
	lookupReceiver := bson.D{{
		"$lookup", bson.D{{"from", PROFILE}, {"localField", "receiver"}, {"foreignField", "id"}, {"as", "receiver"}},
	}}
	project := bson.D{{"$project", bson.M{
		"_id":          0,
		"receiver._id": 0,
	}}}
	sort := bson.D{{
		"$sort", bson.M{"created_at": -1},
	}}

	cursor, err := repo.DB.Collection(CHAT_ROOM).Aggregate(repo.Ctx, mongo.Pipeline{
		match,
		lookupReceiver,
		project,
		sort,
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(repo.Ctx, &chatList)
	if err != nil {
		return nil, util.GetError("error_retrieve")
	}
	return chatList, nil
}

// GetChatCounter ...
func (repo *ChatRepository) GetChatCounter(sender string) (int64, error) {
	filter := bson.M{"sender": sender, "is_read": false}
	chats, err := repo.DB.Collection(CHAT_ROOM).CountDocuments(repo.Ctx, filter, nil)
	if err != nil {
		return 0, err
	}
	return chats, nil
}
