package server

import (
	"encoding/json"
	"errors"
	"log"
	"readygo/wesocket-chat/db"
	"readygo/wesocket-chat/model"
	"readygo/wesocket-chat/payload"
	"sync"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
)

// ServerWs websocket 服务
type ServerWs struct {
	Clients     map[string]*websocket.Conn
	PayloadChan chan payload.Payload
	Repo        *db.Repo
	Mutex       sync.RWMutex
}

// NewServerWs 创建一个ServerWs，初始化成员字段
func NewServerWs(repo *db.Repo) *ServerWs {
	return &ServerWs{
		Clients:     make(map[string]*websocket.Conn),
		PayloadChan: make(chan payload.Payload),
		Repo:        repo,
	}
}

func (srv *ServerWs) Handler() {
	for data := range srv.PayloadChan {
		switch data.Action {
		case payload.RegisterAct:
			srv.Register(&data)
		case payload.LoginAct:
			srv.Login(&data)
		case payload.LogoutAct:
			srv.Logout(&data)
		case payload.BroadcastAct:
			srv.Broadcast(&data)
		case payload.MessageAct:
			fallthrough
		case payload.ErrorAct:
			srv.SendMessage(&data)
		}
	}
}

// AddClient 给ServerWs添加一个client
func (srv *ServerWs) AddClient(conn *websocket.Conn, id string) {
	srv.Mutex.RLock()
	defer srv.Mutex.RUnlock()

	_, exists := srv.Clients[id]
	if exists {
		return
	}

	srv.Clients[id] = conn
}

func (srv *ServerWs) Register(p *payload.Payload) {
	// 不管注册成功与否都要删除临时client
	defer srv.DeleteClient(p.ClientID)

	var errPayload *payload.Payload

	// 将 p.Data map 类型转换成 struct
	var req payload.RegisterRequest
	if ok := srv.bindRequest(p, &req); !ok {
		return
	}

	// 检查用户名是否已存在
	users, err := srv.Repo.UserRepo.List()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Println("检查用户错误: ", err)
		errPayload = payload.NewErrorPayload(p.ClientID, "服务错误，请稍后重试")
		srv.SendMessage(errPayload)
		return
	}

	for _, user := range users {
		if user.Nickname == req.Nickname {
			log.Println("用户已存在， nickname: ", req.Nickname, p.ClientID)
			errPayload = payload.NewErrorPayload(p.ClientID, "用户名已存在，请更换一个")
			srv.SendMessage(errPayload)
			return
		}
	}

	user, err := model.NewUser(req.Nickname, req.Password, "/static/images/default.png")
	if err != nil {
		log.Println("创建用户失败: ", err)
		errPayload = payload.NewErrorPayload(p.ClientID, "创建用户失败")
		srv.SendMessage(errPayload)
		return
	}

	err = srv.Repo.UserRepo.Save(user)
	if err != nil {
		log.Println("保存用户失败： ", err)
		errPayload = payload.NewErrorPayload(p.ClientID, "保存用户失败")
		srv.SendMessage(errPayload)
		return
	}

	// 注册成功
	messagePayload := payload.NewMessagePayload(p.ClientID, "注册成功")

	srv.SendMessage(messagePayload)
}

func (srv *ServerWs) Login(p *payload.Payload) {
	defer srv.DeleteClient(p.ClientID)

	var errPayload *payload.Payload

	// 将 p.Data map 类型转换成 struct
	var req payload.LoginRequest
	if ok := srv.bindRequest(p, &req); !ok {
		return
	}

	user, err := srv.Repo.UserRepo.Get(req.Nickname)
	if err != nil {
		log.Println("查找用户失败 ", err)
		var msg = "查找用户失败"
		if err == redis.Nil {
			msg = "用户不存在"
		}
		errPayload = payload.NewErrorPayload(p.ClientID, msg)
		srv.SendMessage(errPayload)
		return
	}

	if ok := user.CheckedPswd(req.Password, user.Password); !ok {
		log.Println("密码不匹配")
		errPayload = payload.NewErrorPayload(p.ClientID, "密码不匹配")
		srv.SendMessage(errPayload)
		return
	}

	conn := srv.GetConn(p.ClientID)
	if conn == nil {
		return
	}

	srv.AddClient(conn, user.ID.String())

	cookie, err := conn.Request().Cookie("session")
	if err != nil {
		log.Println("cookie error: ", err)
		return
	}
	log.Println("cookie val: ", cookie.Value)

	srv.SendMessage(payload.NewMessagePayload(user.ID.String(), "登录成功"))
}

func (srv *ServerWs) Logout(payload *payload.Payload) {

}

func (srv *ServerWs) SendMessage(payload *payload.Payload) error {
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

func (srv *ServerWs) Broadcast(payload *payload.Payload) {

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

	// conn := srv.GetConn(id)
	// if conn != nil {
	// 	conn.Close()
	// }

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

func (srv *ServerWs) writeCookie() {
	// server := websocket.Server{Handler: websocket.Handler(func(c *websocket.Conn) {})}
	// server.ServeHTTP()
}
