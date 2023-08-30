package model

import (
	"log"
	"time"

	"github.com/google/uuid"
)

// Session 会话
type Session struct {
	ID      uuid.UUID `redis:"id" json:"id"`
	UserID  uuid.UUID `redis:"userId" json:"userId"`
	Expires time.Time `redis:"expires" json:"expires"`
}

func NewSession(userID uuid.UUID) (*Session, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		log.Println("uuid.NewRandom failed: ", err)
		return nil, err
	}

	sess := &Session{
		ID:      id,
		UserID:  userID,
		Expires: time.Now().Add(24 * time.Hour),
	}

	return sess, nil
}
