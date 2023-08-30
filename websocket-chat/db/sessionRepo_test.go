package db_test

import (
	"readygo/wesocket-chat/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func saveSess(t *testing.T) *model.Session {
	user := saveUser(t)
	s := &model.Session{
		ID:      uuid.Must(uuid.NewUUID()),
		UserID:  user.ID,
		Expires: time.Now().Add(time.Hour),
	}
	err := testSessRepo.Save(s)
	require.NoError(t, err)
	return s
}

func TestSaveSess(t *testing.T) {
	saveSess(t)
}

func TestGeteSess(t *testing.T) {
	s := saveSess(t)

	s2, err := testSessRepo.Get(s.ID.String())
	require.NoError(t, err)
	require.Equal(t, s.ID, s2.ID)
	require.Equal(t, s.UserID, s2.UserID)
	require.WithinDuration(t, s.Expires, s2.Expires, time.Second)
}
