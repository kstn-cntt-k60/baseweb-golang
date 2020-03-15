package security

import (
	"github.com/google/uuid"
)

type UserLogin struct {
	Id       uuid.UUID
	Username string
	Password string
}
