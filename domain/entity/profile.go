package entity

import (
	"os"
	"regexp"
	"time"

	"github.com/majid-cj/go-chat-server/util"
)

// AuthorizedProfile ...
var AuthorizedProfile = map[uint8]bool{1: true, 2: false}

const (
	// NICKNAME_PATTERN ...
	NICKNAME_PATTERN = `^(On!.*\.\.)(On!_*__)(On!_.)(On!.*\.$)(On![0-9]*[0-9])[\w.]{2,25}$`
)

// MemberProfile ...
type MemberProfile struct {
	ID           string    `json:"id" bson:"id"`
	Member       string    `json:"member" bson:"member"`
	DisplayName  string    `json:"display_name" bson:"display_name"`
	NickName     string    `json:"nick_name" bson:"nick_name"`
	ProfileImage string    `json:"profile_image" bson:"profile_image"`
	Authorized   uint8     `json:"authorized" bson:"authorized"`
	Private      bool      `json:"is_private" bson:"is_private"`
	CreatedAt    time.Time `json:"created_at" bson:"created_at"`
}

// MemberProfiles ...
type MemberProfiles []MemberProfile

// GetSerializerMemberProfiles ...
func (profiles MemberProfiles) GetSerializerMemberProfiles() []interface{} {
	results := make([]interface{}, len(profiles))
	for index, profile := range profiles {
		results[index] = profile
	}
	return results
}

// ValidateMemberProfile ...
func (profile *MemberProfile) ValidateMemberProfile() error {
	if _, err := regexp.Match(NICKNAME_PATTERN, []byte(profile.NickName)); err != nil {
		return util.GetError("invalid_nickname")
	}
	if len(profile.DisplayName) == 0 || profile.DisplayName == "" {
		return util.GetError("invalid_display_name")
	}
	return nil
}

// PrepareDefaultProfile ...
func (profile *MemberProfile) PrepareDefaultProfile(member, name, email string) {
	profile.ID = util.ULID()
	profile.Member = member
	profile.DisplayName = name
	profile.NickName = util.GenerateNickNameFromEmail(email)
	profile.ProfileImage = os.Getenv("DEFAULT_PROFILE_PIC")
	profile.Private = false
	profile.CreatedAt = util.GetTimeNow()
}
