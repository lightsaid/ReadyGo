package main

import (
	"log"
	"net/http"
	"readygo/wesocket-chat/payload"
	"readygo/wesocket-chat/server"

	"github.com/google/uuid"
	"golang.org/x/net/websocket"
)

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	app.render(w, r, "index.page.html", nil)
}

func (app *application) aboutHandler(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "about.page.html", nil)
}

func (app *application) chatHandler(conn *websocket.Conn, srv *server.ServerWs) {
	defer conn.Close()
	for {
		var req payload.Payload

		err := websocket.JSON.Receive(conn, &req)
		if err != nil {
			// if err == io.EOF {
			// 	log.Println("客户端退出")
			// }
			srv.DeleteClient(req.ClientID)
			// 如果 clientID 没有则是注册或登录业务
			log.Println("->> receive error: ", err, " client ", req.ClientID)
			break
		}

		log.Printf("->> receive: addr: %s, req: %v \n", conn.RemoteAddr(), req)

		go func() {
			defer func() {
				if f := recover(); f != nil {
					log.Printf("recover: %v\n", f)
				}
			}()

			// 为了后续的操作, 生成随机uuid，创建一个临时的客户端
			tempID, err := uuid.NewRandom()
			if err != nil {
				var data = payload.Payload{
					Action: payload.ErrorAct,
					Data:   payload.MessageRequest{Message: "服务错误，请稍后重试"},
				}
				err = websocket.JSON.Send(conn, data)
				if err != nil {
					panic(err)
				}
			}

			// 将 conn 添加到 ServerWs.Clients 方便后续操作
			srv.AddClient(conn, tempID.String())
			req.ClientID = tempID.String()

			srv.PayloadChan <- req
		}()
	}
}
