package types

import (
	"github.com/cloudness-io/cloudness/types/enum"
)

type (
	// User is a user representing an end user.
	User struct {
		// Fields from user
		ID          int64  `db:"principal_id"             json:"-"`
		UID         string `db:"principal_uid"            json:"uid"`
		Email       string `db:"principal_email"          json:"email"`
		DisplayName string `db:"principal_display_name"   json:"display_name"`
		AvatarUrl   string `db:"principal_avatar_url"     json:"avatar_url"`
		Blocked     bool   `db:"principal_blocked"        json:"blocked"`
		Password    string `db:"principal_user_password"  json:"-"`
		Salt        string `db:"principal_salt"           json:"-"`
		Created     int64  `db:"principal_created"        json:"created"`
		Updated     int64  `db:"principal_updated"        json:"updated"`
	}
)

func (u *User) ToPrincipal() *Principal {
	return &Principal{
		ID:          u.ID,
		UID:         u.UID,
		Email:       u.Email,
		Type:        enum.PrincipalTypeUser,
		DisplayName: u.DisplayName,
		Blocked:     u.Blocked,
		Password:    u.Password,
		Salt:        u.Salt,
		Created:     u.Created,
		Updated:     u.Updated,
	}
}

func (u *User) ToPrincipalInfo() *PrincipalInfo {
	return u.ToPrincipal().ToPrincipalInfo()
}
