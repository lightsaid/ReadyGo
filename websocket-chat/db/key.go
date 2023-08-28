package db

const (
	Key_Users    = "users"
	Key_Session  = "session"
	Key_Nickname = "nickname"
)

func Key_User(id string) string {
	return "user:" + id
}
