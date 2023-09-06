package payload

import (
	"readygo/wesocket-chat/model"

	"github.com/google/uuid"
)

// 定义操作类型（或者说协议）
type ActionInt int

const (
	// 发送普通消息
	MessageAct ActionInt = iota + 1
	// 发送广播消息
	BroadcastAct
	// 发送错误消息
	ErrorAct
	// 用户列表
	UsersAct
)

// Payload 定义websocket传参格式
type Payload struct {
	// ServerWs.Clients 的key
	ClientID string      `json:"clientId"`
	Action   ActionInt   `json:"action"`
	Data     interface{} `json:"data"`
}

type UserDto struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Nickname string    `json:"nickname"`
	Avatar   string    `json:"avatar"`
	OnLine   bool      `json:"onLine"`
}

func NewUserDto(user *model.User, online bool) UserDto {
	return UserDto{
		ID:       user.ID,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
		OnLine:   online,
	}
}

// 下面是每个 ActionInt 对应的Data

// MessageRequest 通用格式，适用SendMessageAct、BroadcastAct、ErrorMessageAct
type MessageRequest struct {
	Message string `json:"message"`
}

type BroadcastMessage struct {
	Message string `json:"message"`
	FromID  string `json:"fromId"` // 来自谁？
}

func NewMessagePayload(clientID string, msg string) *Payload {
	return &Payload{
		Action:   MessageAct,
		ClientID: clientID,
		Data: MessageRequest{
			Message: msg,
		},
	}
}

func NewErrorPayload(clientID string, msg string) *Payload {
	return &Payload{
		Action:   ErrorAct,
		ClientID: clientID,
		Data: MessageRequest{
			Message: msg,
		},
	}
}

func NewBroadcastPayload(clientID string, fromId string, msg string) *Payload {
	return &Payload{
		Action:   BroadcastAct,
		ClientID: clientID,
		Data: BroadcastMessage{
			Message: msg,
			FromID:  fromId,
		},
	}
}
