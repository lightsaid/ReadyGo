package db_test

import (
	"log"
	"readygo/wesocket-chat/model"
	"readygo/wesocket-chat/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func saveUser(t *testing.T) *model.User {
	nickname := util.RandomString(4)
	avatar := "http://" + util.RandomString(6)
	user, err := model.NewUser(nickname, "123456", avatar)
	require.NoError(t, err)

	err = testUserRepo.Save(user)
	require.NoError(t, err)

	return user
}

func TestSaveUser(t *testing.T) {
	saveUser(t)
}

func TestGetUser(t *testing.T) {
	user := saveUser(t)

	user2, err := testUserRepo.Get(user.ID.String())
	require.NoError(t, err)

	require.Equal(t, user.ID, user2.ID)
	require.Equal(t, user.Nickname, user2.Nickname)
	require.Equal(t, user.Avatar, user2.Avatar)
}

func TestListUser(t *testing.T) {
	var n = 5
	for i := 0; i < n; i++ {
		saveUser(t)
	}

	users, err := testUserRepo.List()
	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.True(t, len(users) > n-1)
	log.Println("->> ", len(users))
}
