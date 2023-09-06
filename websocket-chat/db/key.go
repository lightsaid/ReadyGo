package db

const (
	Key_Users    = "users"
	Key_Session  = "session"
	Key_Nickname = "nickname"
)

func Key_User(name string) string {
	return "user:" + name
}

func Key_UserID(id string) string {
	// TODO: 检查用户名不能以 “id:” 开头
	return "id:" + id
}
