package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"readygo/wesocket-chat/db"
	"readygo/wesocket-chat/model"
	"readygo/wesocket-chat/payload"
	"readygo/wesocket-chat/server"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
)

type packet map[string]interface{}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("json encode error: ", err)
	}
}

func (app *application) indexHandler(w http.ResponseWriter, r *http.Request, user *model.User) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	app.render(w, r, "index.page.html", user)
}

func (app *application) aboutHandler(w http.ResponseWriter, r *http.Request, uuser *model.User) {
	app.render(w, r, "about.page.html", nil)
}

func (app *application) registerHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}
	var respErr error
	defer func(w http.ResponseWriter, r *http.Request) {
		if respErr != nil {
			resp := packet{"error": respErr.Error()}
			writeJSON(w, r, http.StatusBadRequest, resp)
		}
	}(w, r)

	r.ParseForm()
	nickname := r.Form.Get("nickname")
	password := r.Form.Get("password")

	// 做一下简单校验
	if nickname == "" || password == "" {
		respErr = errors.New("参数必填")
		return
	}

	// 获取所有并检查用户名是否已存在
	users, err := db.RepoObj.UserRepo.List()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Println("checked user error: ", err)
		respErr = errors.New("服务错误，请稍后重试")
		return
	}

	for _, user := range users {
		if user.Nickname == nickname {
			log.Println("用户已存在， nickname: ", nickname)
			respErr = errors.New("用户名已存在，请更换一个")
			return
		}
	}

	respErr = errors.New("创建用户失败")
	user, err := model.NewUser(nickname, password, "/static/images/default.png")
	if err != nil {
		log.Println("new user error: ", err)
		return
	}

	err = db.RepoObj.UserRepo.Save(user)
	if err != nil {
		log.Println("save user error: ", err)
		return
	}

	// 成功操作
	respErr = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	if strings.ToUpper(r.Method) != "POST" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
		return
	}

	var err error
	defer func(w http.ResponseWriter, r *http.Request) {
		if err != nil {
			resp := packet{"error": err.Error()}
			writeJSON(w, r, http.StatusBadRequest, resp)
		}
	}(w, r)

	r.ParseForm()
	nickname := r.Form.Get("nickname")
	password := r.Form.Get("password")

	// 做一下简单校验
	if nickname == "" || password == "" {
		err = errors.New("参数必填")
		return
	}

	// 查找用户
	var user *model.User
	user, err = db.RepoObj.UserRepo.Get(nickname)
	if err != nil {
		log.Println("get user error: ", err)
		if errors.Is(err, redis.Nil) {
			err = errors.New("用户不存在")
		} else {
			err = errors.New("服务内存错误，稍后重试")
		}
		return
	}

	// 检查密码
	if ok := user.CheckedPswd(password, user.Password); !ok {
		log.Println("密码不匹配")
		err = errors.New("密码不匹配")
		return
	}

	// 写入 session
	var ss *model.Session
	ss, err = model.NewSession(user.ID)
	if err != nil {
		log.Println("new session error: ", err)
		err = errors.New("登录失败")
		return
	}
	err = db.RepoObj.SessionRepo.Save(ss)
	if err != nil {
		log.Println("save session error: ", err)
		err = errors.New("登录失败")
		return
	}

	// 写入cookie
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: ss.ID.String(),
	})

	// 成功操作
	err = nil
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) wsHandler(conn *websocket.Conn) {
	defer conn.Close()
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
			app.serverWs.AddClient(conn, tempID.String())
			req.ClientID = tempID.String()

			app.serverWs.PayloadChan <- req
		}()
	}
}

func (app *application) wsMiddleware(w http.ResponseWriter, r *http.Request) {
	wssrv := websocket.Server{
		Handler: websocket.Handler(app.wsHandler),
	}

	cookieVal, ok := r.Context().Value("cookie").(string)
	if ok {
		fmt.Println("->>> cookieVal: ", cookieVal)
	}

	fmt.Println("->> ws middleware")

	wssrv.ServeHTTP(w, r)
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
