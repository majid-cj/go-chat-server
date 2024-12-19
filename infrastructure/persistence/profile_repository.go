package persistence

import (
	"context"
	"strings"

	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/domain/repository"
	"github.com/majid-cj/go-chat-server/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MemberProfileRepository ...
type MemberProfileRepository struct {
	Ctx context.Context
	DB  *mongo.Collection
}

// NewMemberProfileRepository ...
func NewMemberProfileRepository(db *mongo.Database) *MemberProfileRepository {
	return &MemberProfileRepository{
		Ctx: context.Background(),
		DB:  db.Collection(PROFILE),
	}
}

var _ repository.ProfileRepository = &MemberProfileRepository{}

// CreateMemberProfile ...
func (repo *MemberProfileRepository) CreateMemberProfile(profile *entity.MemberProfile) (*entity.MemberProfile, error) {
	_, err := repo.DB.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.MDoc{"nick_name": bsonx.Int64(1)},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, util.GetError("general_error")
	}

	_, err = repo.DB.InsertOne(repo.Ctx, profile)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, util.GetError("nickname_taken")
		}
		return nil, util.GetError("general_error")
	}
	return profile, nil
}

// UpdateMemberProfile ...
func (repo *MemberProfileRepository) UpdateMemberProfile(profile *entity.MemberProfile) (*entity.MemberProfile, error) {
	filter := bson.M{"id": profile.ID}
	update := bson.M{"$set": profile}
	_, err := repo.DB.UpdateOne(repo.Ctx, filter, update)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, util.GetError("nickname_taken")
		}
		return nil, util.GetError("general_error")
	}
	return profile, nil
}

// GetMemberProfileByID ...
func (repo *MemberProfileRepository) GetMemberProfileByID(ID string) (*entity.MemberProfile, error) {
	var profile entity.MemberProfile
	filter := bson.M{"id": ID}
	err := repo.DB.FindOne(repo.Ctx, filter).Decode(&profile)
	if err != nil {
		return nil, util.GetError("profile_not_found")
	}
	return &profile, nil
}

// GetMemberProfileByMemberID ...
func (repo *MemberProfileRepository) GetMemberProfileByMemberID(member string) (*entity.MemberProfile, error) {
	var profile entity.MemberProfile
	filter := bson.M{"member": member}
	err := repo.DB.FindOne(repo.Ctx, filter).Decode(&profile)
	if err != nil {
		return nil, util.GetError("profile_not_found")
	}
	return &profile, nil
}

// GetMemberProfileByNickName ...
func (repo *MemberProfileRepository) GetMemberProfileByNickName(nickname string) (*entity.MemberProfile, error) {
	var profile entity.MemberProfile

	search := bson.M{
		"nick_name": bsonx.Regex(nickname, "i"),
	}

	err := repo.DB.FindOne(repo.Ctx, search).Decode(&profile)
	if err != nil {
		return nil, util.GetError("profile_not_found")
	}

	return &profile, nil
}
