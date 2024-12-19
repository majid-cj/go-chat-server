package persistence

import (
	"context"

	"github.com/majid-cj/go-chat-server/domain/entity"
	"github.com/majid-cj/go-chat-server/util"
	"github.com/majid-cj/go-chat-server/util/security"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// VerifyCodeRepository ...
type VerifyCodeRepository struct {
	Ctx      context.Context
	Db       *mongo.Collection
	DbMember *mongo.Collection
}

// NewVerifyCodeRepository ...
func NewVerifyCodeRepository(db *mongo.Database) *VerifyCodeRepository {
	return &VerifyCodeRepository{
		Ctx:      context.Background(),
		Db:       db.Collection(VERIFY_CODE),
		DbMember: db.Collection(MEMBER),
	}
}

// CreateVerificationCode ...
func (repo *VerifyCodeRepository) CreateVerificationCode(code *entity.VerificationCode) (*entity.VerificationCode, error) {
	filter := bson.M{"member": code.Member}
	repo.Db.DeleteMany(repo.Ctx, filter)
	_, err := repo.Db.InsertOne(repo.Ctx, code)
	if err != nil {
		return nil, err
	}
	return code, nil
}

// CreateVerificationCodeFromEmail ...
func (repo *VerifyCodeRepository) CreateVerificationCodeFromEmail(code *entity.VerificationCode) (*entity.VerificationCode, error) {
	var member entity.Member
	var verifyCode entity.VerificationCode
	err := repo.DbMember.FindOne(repo.Ctx, bson.M{"email": code.Email, "member_type": 3, "source": 1}).Decode(&member)
	if err != nil {
		return nil, util.GetError("no_email_account")
	}

	repo.Db.DeleteMany(repo.Ctx, bson.M{"member": member.ID})
	verifyCode.PrepareVerificationCode(member.ID, code.CodeType)
	_, err = repo.Db.InsertOne(repo.Ctx, &verifyCode)
	if err != nil {
		return nil, util.GetError("general_error")
	}
	return &verifyCode, nil
}

// ResetPassword ...
func (repo *VerifyCodeRepository) ResetPassword(code *entity.VerificationCode) error {
	var verifyCode entity.VerificationCode
	var member entity.Member
	filterMember := bson.M{"email": code.Email, "member_type": 3, "source": 1}
	err := repo.DbMember.FindOne(repo.Ctx, filterMember).Decode(&member)
	if err != nil {
		return util.GetError("general_error")
	}

	filterCode := bson.M{"member": member.ID, "code": code.Code, "code_type": code.CodeType, "taken": false}
	err = repo.Db.FindOne(repo.Ctx, filterCode).Decode(&verifyCode)
	if err != nil {
		return util.GetError("general_error")
	}

	if util.GetTimeNow().Unix() > verifyCode.ExpiredAt.Unix() {
		return util.GetError("token_expired")
	}

	saltedPassword := security.NewPassHash(member.ID, member.Email, code.Password, []byte{128})
	_, err = repo.Db.UpdateMany(repo.Ctx, filterCode, bson.M{"$set": bson.M{"taken": true}})
	if err != nil {
		return util.GetError("general_error")
	}

	member.Password = saltedPassword
	member.UpdateAt = util.GetTimeNow()
	_, err = repo.DbMember.UpdateOne(repo.Ctx, bson.M{"id": member.ID}, bson.M{"$set": member})
	if err != nil {
		return util.GetError("general_error")
	}

	return nil
}

// CheckVerificationCode ...
func (repo *VerifyCodeRepository) CheckVerificationCode(code *entity.VerificationCode) error {
	var verifyCode entity.VerificationCode
	filter := bson.M{"member": code.Member, "code": code.Code, "code_type": code.CodeType}

	err := repo.Db.FindOne(repo.Ctx, filter).Decode(&verifyCode)
	if err != nil {
		return util.GetError("general_error")
	}

	if util.GetTimeNow().Unix() > verifyCode.ExpiredAt.Unix() {
		return util.GetError("token_expired")
	}

	update := bson.M{"$set": bson.M{"taken": true}}
	_, err = repo.Db.UpdateOne(repo.Ctx, filter, update)
	if err != nil {
		return util.GetError("general_error")
	}

	_, err = repo.DbMember.UpdateOne(repo.Ctx, bson.M{"id": code.Member}, bson.M{"$set": bson.M{"verified": true}})
	if err != nil {
		return util.GetError("general_error")
	}
	return nil
}

// RenewVerificationCode ...
func (repo *VerifyCodeRepository) RenewVerificationCode(code *entity.VerificationCode) (*entity.VerificationCode, error) {
	var verifyCode entity.VerificationCode
	filter := bson.M{"member": code.Member}
	_, err := repo.Db.DeleteMany(repo.Ctx, filter)
	if err != nil {
		return nil, util.GetError("general_error")
	}
	verifyCode.PrepareVerificationCode(code.Member, code.CodeType)
	_, err = repo.Db.InsertOne(repo.Ctx, verifyCode)
	if err != nil {
		return nil, util.GetError("general_error")
	}
	return &verifyCode, nil
}
