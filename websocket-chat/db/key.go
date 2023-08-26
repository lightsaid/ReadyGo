package db

import "strconv"

const (
	Key_Users    = "users"
	Key_Session  = "session"
	Key_Nickname = "nickname"
)

func Key_User(id int) string {
	return "user:" + strconv.Itoa(id)
}
