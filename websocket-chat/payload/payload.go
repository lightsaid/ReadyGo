package payload

// 定义操作类型（或者说协议）
type ActionInt int

const (
	// 注册
	RegisterAct ActionInt = iota + 1
	// 登录
	LoginAct
	// 退出
	LogoutAct
	// 发送消息
	MessageAct
	// 广播
	BroadcastAct
	// 发送错误消息
	ErrorAct
)

// Payload 定义websocket传参格式
type Payload struct {
	// ServerWs.Clients 的key
	ClientID string      `json:"clientId"`
	Action   ActionInt   `json:"action"`
	Data     interface{} `json:"data"`
}

// 下面是每个 ActionInt 对应的Data

type RegisterRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	UserID string `json:"userId"`
}

// MessageRequest 通用格式，适用SendMessageAct、BroadcastAct、ErrorMessageAct
type MessageRequest struct {
	Message string `json:"message"`
}

func NewMessagePayload(clientID string, msg string) Payload {
	return Payload{
		Action:   MessageAct,
		ClientID: clientID,
		Data: MessageRequest{
			Message: msg,
		},
	}
}

func NewErrorPayload(clientID string, msg string) Payload {
	return Payload{
		Action:   ErrorAct,
		ClientID: clientID,
		Data: MessageRequest{
			Message: msg,
		},
	}
}
