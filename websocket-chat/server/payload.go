package server

import "readygo/wesocket-chat/model"

type ActionType string

const (
	RegisterOpt ActionType = "REGISTER"
	LoginOpt    ActionType = "LOGIN"
	LogoutOpt   ActionType = "LOGOUT"
	Broadcast   ActionType = "BROADCAST"
)

type Payload struct {
	Action ActionType  `json:"action"`
	Data   interface{} `json:"data"`
}

// 下面是每个 ActionType 对应的Data

type RegisterRequest struct {
	User model.User `json:"user"`
}

type LoginOptRequest struct {
	User model.User `json:"user"`
}

type LogoutOptRequest struct {
	UserID int `json:"userId"`
}

type BroadcastRequest struct {
	Message string `json:"message"`
}
