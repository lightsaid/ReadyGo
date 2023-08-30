package main

import (
	"net/http"
	"readygo/wesocket-chat/db"
	"readygo/wesocket-chat/server"

	"github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
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

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		app.wsMiddleware(w, r)
	})

	mux.Handle("/chat", websocket.Handler(func(c *websocket.Conn) {
		app.chatHandler(c, app.serverWs)
	}))

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}

func (app *application) myhandler(w http.ResponseWriter, r *http.Request) http.Handler {
	return websocket.Handler(func(c *websocket.Conn) {

	})
}
