package entity

import (
	"time"

	"github.com/majid-cj/go-chat-server/util"
	"github.com/majid-cj/go-chat-server/util/security"
)

// Member ...
type Member struct {
	ID        string            `bson:"id" json:"id"`
	Email     string            `bson:"email" json:"email"`
	Password  security.PassHash `bson:"password" json:"password"`
	Verified  bool              `bson:"verified" json:"verified"`
	Active    bool              `bson:"active" json:"active"`
	CreatedAt time.Time         `bson:"created_at"`
	UpdateAt  time.Time         `bson:"update_at"`
}

// MemberSerializer ...
type MemberSerializer struct {
	ID        string `json:"id"`
	Type      uint8  `json:"member_type"`
	Email     string `json:"email"`
	Verified  bool   `json:"verified"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
}

// Members ...
type Members []Member

// PrepareMember ...
func (m *Member) PrepareMember() {
	m.ID = util.ULID()
	m.Email = util.EscapeString(m.Email)
	m.Active = true
	m.Verified = false
	m.CreatedAt = util.GetTimeNow()
	m.UpdateAt = util.GetTimeNow()
}

// PrepareSocialMember ...
func (m *Member) PrepareSocialMember() {
	m.ID = util.ULID()
	m.Email = util.EscapeString(m.Email)
	m.Active = true
	m.Verified = true
	m.CreatedAt = util.GetTimeNow()
	m.UpdateAt = util.GetTimeNow()
}

// GetMemberSerializer ...
func (m Member) GetMemberSerializer() MemberSerializer {
	return MemberSerializer{
		ID:        m.ID,
		Email:     m.Email,
		Verified:  m.Verified,
		Active:    m.Active,
		CreatedAt: m.CreatedAt.String(),
	}
}

// GetMembersSerializer ...
func (members Members) GetMembersSerializer() []interface{} {
	results := make([]interface{}, len(members))
	for index, user := range members {
		results[index] = user.GetMemberSerializer()
	}
	return results
}
