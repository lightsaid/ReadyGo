package server

import (
	"encoding/json"
	"fmt"
	"log"
	"readygo/wesocket-chat/db"
	"readygo/wesocket-chat/payload"
	"sync"

	"golang.org/x/net/websocket"
)

// ServerWs websocket 服务
type ServerWs struct {
	Clients     map[string]*websocket.Conn
	PayloadChan chan *payload.Payload
	Repo        *db.Repo
	Mutex       sync.RWMutex
}

// NewServerWs 创建一个ServerWs，初始化成员字段
func NewServerWs(repo *db.Repo) *ServerWs {
	return &ServerWs{
		Clients:     make(map[string]*websocket.Conn),
		PayloadChan: make(chan *payload.Payload),
		Repo:        repo,
	}
}

func (srv *ServerWs) Handler() {
	for data := range srv.PayloadChan {
		switch data.Action {
		case payload.BroadcastAct:
			srv.Broadcast(data)
		case payload.UsersAct:
			srv.AllUsers(data)
		case payload.MessageAct:
			fallthrough
		case payload.ErrorAct:
			srv.SendMessage(data)
		}
	}
}

// AddClient 给ServerWs添加一个client
func (srv *ServerWs) AddClient(conn *websocket.Conn, id string) {
	fmt.Println("add client ", id)
	srv.Mutex.RLock()
	defer srv.Mutex.RUnlock()

	srv.Clients[id] = conn
}

func (srv *ServerWs) AllUsers(p *payload.Payload) {
	users, err := srv.Repo.UserRepo.List()
	if err != nil {
		log.Println("AllUsers() error: ", err)
		errPayload := payload.NewErrorPayload(p.ClientID, "获取用列表错误")
		srv.SendMessage(errPayload)
		return
	}

	var res []payload.UserDto
	for _, user := range users {
		_, exists := srv.Clients[user.ID.String()]
		res = append(res, payload.NewUserDto(user, exists))
	}

	p.Data = res

	srv.SendMessage(p)
}

func (srv *ServerWs) Broadcast(p *payload.Payload) {
	srv.Mutex.Lock()
	defer srv.Mutex.Unlock()

	var req payload.BroadcastMessage
	if okay := srv.bindRequest(p, &req); !okay {
		return
	}

	for id := range srv.Clients {
		data := payload.NewBroadcastPayload(id, req.FromID, req.Message)
		srv.SendMessage(data)
	}
}

func (srv *ServerWs) SendMessage(payload *payload.Payload) error {
	// TODO: 速率控制
	conn, exists := srv.Clients[payload.ClientID]
	if exists {
		err := websocket.JSON.Send(conn, payload)
		if err != nil {
			log.Println("SendMessage error: ", err)
			return err
		}
	} else {
		log.Println("client 不存在: ", payload.ClientID, payload.Data)
	}
	return nil
}

func (srv *ServerWs) OnLine() int {
	return len(srv.Clients)
}

func (srv *ServerWs) GetConn(id string) *websocket.Conn {
	return srv.Clients[id]
}

func (srv *ServerWs) DeleteClient(id string) {
	srv.Mutex.RLock()
	defer srv.Mutex.RUnlock()

	delete(srv.Clients, id)
}

func (srv *ServerWs) bindRequest(p *payload.Payload, req interface{}) bool {
	var errPayload *payload.Payload

	buf, err := json.Marshal(p.Data)
	if err != nil {
		log.Println("register json.Marshal error: ", err)
		errPayload = payload.NewErrorPayload(p.ClientID, "参数错误")
		srv.SendMessage(errPayload)
		return false
	}

	err = json.Unmarshal(buf, &req)
	if err != nil {
		log.Println("register json.Unmarshal error: ", err)
		errPayload = payload.NewErrorPayload(p.ClientID, "参数错误")
		srv.SendMessage(errPayload)
		return false
	}

	return true
}
