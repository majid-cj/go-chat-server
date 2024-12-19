package repository

import "github.com/majid-cj/go-chat-server/domain/entity"

// ProfileRepository ...
type ProfileRepository interface {
	CreateMemberProfile(*entity.MemberProfile) (*entity.MemberProfile, error)
	UpdateMemberProfile(*entity.MemberProfile) (*entity.MemberProfile, error)
	GetMemberProfileByID(string) (*entity.MemberProfile, error)
	GetMemberProfileByMemberID(string) (*entity.MemberProfile, error)
	GetMemberProfileByNickName(string) (*entity.MemberProfile, error)
}
