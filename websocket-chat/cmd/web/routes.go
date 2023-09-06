package main

import (
	"fmt"
	"net/http"
	"readygo/wesocket-chat/db"
	"readygo/wesocket-chat/model"
	"readygo/wesocket-chat/server"

	"github.com/redis/go-redis/v9"
)

func (app *application) routes(rdb *redis.Client) http.Handler {
	mux := http.NewServeMux()

	store := db.NewRepo(rdb)

	app.serverWs = server.NewServerWs(store)

	go app.serverWs.Handler()

	mux.Handle("/", app.middlewareAuth(app.indexHandler))
	mux.Handle("/about", app.middlewareAuth(app.aboutHandler))

	mux.HandleFunc("/register", app.registerHandler)
	mux.HandleFunc("/login", app.loginHandler)

	mux.HandleFunc("/logout", app.logoutHandler)

	// websocket 路由
	mux.Handle("/ws", app.middlewareAuth(func(w http.ResponseWriter, r *http.Request, u *model.User) {
		fmt.Println("ws ->> ", "app.middlewareAuth")
		app.wsMiddleware(w, r, u)
	}))

	// 静态资源
	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return app.recoverPanic(mux)
}
