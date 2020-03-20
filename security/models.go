package security

import (
	"time"

	"github.com/google/uuid"
)

type UserLogin struct {
	Id       uuid.UUID
	Username string
	Password string
}

type ClientUserLogin struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type Permission struct {
	Id        int16     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type Group struct {
	Id        int16     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

type GroupPermission struct {
	Id struct {
		GroupId      int16 `json:"securityGroupId"`
		PermissionId int16 `json:"securityPermissionId"`
	} `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
