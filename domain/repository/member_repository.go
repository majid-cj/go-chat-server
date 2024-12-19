package repository

import "github.com/majid-cj/go-chat-server/domain/entity"

// MemberRepository ...
type MemberRepository interface {
	CreateMember(*entity.Member) (*entity.Member, error)
	DeleteMember(string) error
	GetMembers() ([]entity.Member, error)
	GetMember(string) (*entity.Member, error)
	GetMembersByType(uint8) ([]entity.Member, error)
	GetMembersBySource(uint8) ([]entity.Member, error)
	GetMemberByEmailAndPassword(*entity.SignUp) (*entity.Member, error)
	GetMemberByEmailAndSource(*entity.Member) (*entity.Member, uint8, error)
	UpdatePassword(*entity.Member) error
}
