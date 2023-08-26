package db_test

import (
	"readygo/wesocket-chat/model"
	"readygo/wesocket-chat/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func save(t *testing.T) model.User {
	user := model.User{
		ID:       int(util.RandomInt(1000, 100000)),
		Nickname: util.RandomString(4),
		Avatar:   "http://",
		Password: "123456",
	}
	err := user.GenHashedPswd()
	require.NoError(t, err)

	err = testUserRepo.Save(user)
	require.NoError(t, err)

	user2, err := testUserRepo.Get(user.ID)
	require.NoError(t, err)

	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Nickname, user2.Nickname)
	require.Equal(t, user.Avatar, user2.Avatar)

	return user
}

func TestSave(t *testing.T) {
	save(t)
}
