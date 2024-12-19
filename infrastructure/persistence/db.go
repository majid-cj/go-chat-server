package persistence

import (
	"context"
	"fmt"
	"os"

	"github.com/majid-cj/go-chat-server/domain/repository"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository ...
type Repository struct {
	Member     repository.MemberRepository
	VerifyCode repository.VerificationCodeRepository
	Profile    repository.ProfileRepository
	Chat       repository.ChatRepository
	Ctx        context.Context
	Client     *mongo.Client
}

// NewRepository ...
func NewRepository() (*Repository, error) {
	URL := fmt.Sprintf("%s://%s:%s@%s/%s", os.Getenv("DB_DRIVER"), os.Getenv("DB_HOST"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	clientOption := options.Client().ApplyURI(URL).SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1))
	ctx := context.Background()

	client, err := mongo.Connect(ctx, clientOption)

	if err != nil {
		return nil, err
	}
	db := client.Database(os.Getenv("DB_NAME"))
	return &Repository{
		Member:     NewMemberRepository(db),
		VerifyCode: NewVerifyCodeRepository(db),
		Profile:    NewMemberProfileRepository(db),
		Chat:       NewChatRepository(db),
		Ctx:        ctx,
		Client:     client,
	}, nil
}

const (
	// MEMBER ...
	MEMBER = "member"
	// VERIFY_CODE ...
	VERIFY_CODE = "verify_code"
	// PROFILE ...
	PROFILE = "profile"
	// CHAT ...
	CHAT = "chat"
	// CHAT_ROOM ...
	CHAT_ROOM = "chat_room"
)
