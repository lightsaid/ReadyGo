package main

import (
	"fmt"
	"log"
	"net/http"
	"readygo/wesocket-chat/model"
	"readygo/wesocket-chat/payload"

	"golang.org/x/net/websocket"
)

// wsHandler 负责处理websocket工作
func (app *application) wsHandler(conn *websocket.Conn) {
	defer conn.Close()

	// 第一次进入做一些初始化操作
	go func() {
		defer func() {
			if f := recover(); f != nil {
				log.Printf("recover: %v\n", f)
			}
		}()

		user := app.contextGetUser(conn.Request())
		if user == nil {
			log.Println("建立连接时，上下文获取用户信息失败")
			websocket.JSON.Send(conn, payload.NewErrorPayload("", "服务内部错误，获取用户信息失败"))
			return
		}

		// 将 conn 添加到 ServerWs.Clients 方便后续操作
		app.serverWs.AddClient(conn, user.ID.String())

		msgs := payload.NewMessagePayload(user.ID.String(), "第一次发送消息, websocket 建立连接成功")
		app.serverWs.PayloadChan <- msgs
	}()

	for {
		var req payload.Payload
		err := websocket.JSON.Receive(conn, &req)

		if err != nil {
			// if err == io.EOF {
			// 	log.Println("客户端退出")
			// }
			app.serverWs.DeleteClient(req.ClientID)
			// 如果 clientID 没有则是注册或登录业务
			log.Println("->> receive error: ", err, " client ", req.ClientID)
			break
		}

		fmt.Println("请求入参：", req.ClientID, " - ", req.Action, " - ", req.Data)
		app.serverWs.PayloadChan <- &req
	}
}

func (app *application) wsMiddleware(w http.ResponseWriter, r *http.Request, user *model.User) {
	// 传递 websocket 处理函数 创建一个 ws server
	wssrv := websocket.Server{
		Handler: websocket.Handler(app.wsHandler),
	}

	if user == nil {
		fmt.Println("服务错误，建立ws时获取用户信息失败")
		return
	}

	// 设置 user 上下文
	r = app.contextSetUser(r, user)

	fmt.Println("->> ws middleware ", user.Nickname)

	// 执行 websocket.Server 的 ServeHTTP，同时将 http 请求交由websocker处理
	wssrv.ServeHTTP(w, r)
}
