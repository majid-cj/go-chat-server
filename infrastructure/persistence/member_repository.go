package persistence

import (
	"context"
	"strings"

	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/domain/repository"
	"github.com/majid-cj/go-chat-server/util"
	"github.com/majid-cj/go-chat-server/util/security"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// MemberRepository ...
type MemberRepository struct {
	Ctx context.Context
	DB  *mongo.Collection
}

// NewMemberRepository ...
func NewMemberRepository(db *mongo.Database) *MemberRepository {
	return &MemberRepository{
		Ctx: context.Background(),
		DB:  db.Collection(MEMBER),
	}
}

var _ repository.MemberRepository = &MemberRepository{}

// CreateMember ...
func (repo *MemberRepository) CreateMember(member *entity.Member) (*entity.Member, error) {
	_, err := repo.DB.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys:    bsonx.MDoc{"email": bsonx.Int64(1)},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		return nil, util.GetError("general_error")
	}

	_, err = repo.DB.InsertOne(repo.Ctx, member)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return nil, util.GetError("email_taken")
		} else {
			return nil, util.GetError("general_error")
		}
	}
	return member, nil
}

// DeleteMember ...
func (repo *MemberRepository) DeleteMember(ID string) error {
	filter := bson.M{"id": ID}
	_, err := repo.DB.DeleteOne(repo.Ctx, filter)
	if err != nil {
		return util.GetError("general_error")
	}
	return nil
}

// GetMembers ...
func (repo *MemberRepository) GetMembers() ([]entity.Member, error) {
	var members entity.Members
	filter := bson.M{}
	cursor, err := repo.DB.Find(repo.Ctx, filter, nil)
	if err != nil {
		return nil, util.GetError("general_error")
	}

	defer cursor.Close(repo.Ctx)

	err = cursor.All(repo.Ctx, &members)
	if err != nil {
		return nil, util.GetError("error_retrieve")
	}

	if len(members) == 0 {
		return nil, util.GetError("empty_list")
	}
	return members, nil
}

// GetMember ...
func (repo *MemberRepository) GetMember(ID string) (*entity.Member, error) {
	var member entity.Member
	filter := bson.M{"id": ID}
	err := repo.DB.FindOne(repo.Ctx, filter).Decode(&member)
	if err != nil {
		return nil, util.GetError("member_not_found")
	}
	return &member, nil
}

// GetMembersByType ...
func (repo *MemberRepository) GetMembersByType(memberType uint8) ([]entity.Member, error) {
	var members entity.Members

	filter := bson.M{"member_type": memberType}

	cursor, err := repo.DB.Find(repo.Ctx, filter, nil)
	if err != nil {
		return nil, util.GetError("general_error")
	}

	defer cursor.Close(repo.Ctx)

	for cursor.Next(repo.Ctx) {
		var member entity.Member
		err := cursor.Decode(member)
		if err != nil {
			return nil, util.GetError("general_error")
		}
		members = append(members, member)
	}

	if len(members) == 0 {
		return nil, util.GetError("empty_list")
	}
	return members, nil
}

// GetMembersBySource ...
func (repo *MemberRepository) GetMembersBySource(source uint8) ([]entity.Member, error) {
	var members entity.Members

	filter := bson.M{"source": source}

	cursor, err := repo.DB.Find(repo.Ctx, filter, nil)
	if err != nil {
		return nil, util.GetError("general_error")
	}

	defer cursor.Close(repo.Ctx)

	for cursor.Next(repo.Ctx) {
		var member entity.Member
		err := cursor.Decode(member)
		if err != nil {
			return nil, util.GetError("general_error")
		}
		members = append(members, member)
	}

	if len(members) == 0 {
		return nil, util.GetError("empty_list")
	}
	return members, nil
}

// GetMemberByEmailAndPassword ...
func (repo *MemberRepository) GetMemberByEmailAndPassword(member *entity.SignUp) (*entity.Member, error) {
	var getMember entity.Member

	filter := bson.M{
		"email": member.Email,
	}

	err := repo.DB.FindOne(repo.Ctx, filter).Decode(&getMember)
	if err != nil {
		return nil, util.GetError("email_password_wrong")
	}

	valid := security.EqualPassHash(getMember.ID, member.Email, member.Password, getMember.Password)
	if !valid {
		return nil, util.GetError("email_password_wrong")
	}

	return &getMember, nil
}

// GetMemberByEmailAndSource ...
func (repo *MemberRepository) GetMemberByEmailAndSource(member *entity.Member) (*entity.Member, uint8, error) {
	var getmember entity.Member

	filter := bson.M{
		"email": member.Email,
	}

	err := repo.DB.FindOne(repo.Ctx, filter).Decode(&getmember)

	if err == nil {
		return &getmember, 0, nil
	}

	member.PrepareSocialMember()
	getnewmember, getError := repo.CreateMember(member)

	if getError != nil {
		return &getmember, 0, util.GetError("member_not_found")
	}
	return getnewmember, 1, nil
}

// UpdatePassword ...
func (repo *MemberRepository) UpdatePassword(member *entity.Member) error {
	filter := bson.M{"id": member.ID}
	update := bson.M{"$set": bson.M{
		"password":  member.Password,
		"update_at": member.UpdateAt,
	}}
	_, err := repo.DB.UpdateOne(repo.Ctx, filter, update)
	if err != nil {
		return util.GetError("general_error")
	}
	return nil
}
