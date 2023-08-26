package model

import (
	"time"

	"github.com/google/uuid"
)

// Session 会话
type Session struct {
	ID      uuid.UUID `redis:"id" json:"id"`
	UserID  int       `redis:"userId" json:"userId"`
	Expires time.Time `redis:"expires" json:"expires"`
}
