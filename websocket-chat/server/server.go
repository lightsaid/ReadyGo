package server

import (
	"golang.org/x/net/websocket"
)

// ServerWs websocket 服务
type ServerWs struct {
	Clients     map[int]websocket.Conn
	PayloadChan chan Payload
}

// NewServerWs 创建一个ServerWs，初始化成员字段
func NewServerWs() *ServerWs {
	return &ServerWs{
		Clients:     make(map[int]websocket.Conn),
		PayloadChan: make(chan Payload),
	}
}

func (srv *ServerWs) Handler() {
	for payload := range srv.PayloadChan {
		switch payload.Action {
		case RegisterOpt:
			srv.Register(&payload)
		case LoginOpt:
			srv.Login(&payload)
		case LogoutOpt:
			srv.Logout(&payload)
		case Broadcast:
			srv.Broadcast(&payload)
		}
	}

}

func (srv *ServerWs) Register(payload *Payload) {

}

func (srv *ServerWs) Login(payload *Payload) {

}

func (srv *ServerWs) Logout(payload *Payload) {

}

func (srv *ServerWs) SendMessage(payload *Payload) {

}

func (srv *ServerWs) Broadcast(payload *Payload) {

}

func (srv *ServerWs) OnLine() int {
	return len(srv.Clients)
}
