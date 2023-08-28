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

	srv := server.NewServerWs(store)

	go srv.Handler()

	mux.HandleFunc("/", app.indexHandler)
	mux.HandleFunc("/about", app.aboutHandler)

	mux.Handle("/chat", websocket.Handler(func(c *websocket.Conn) {
		app.chatHandler(c, srv)
	}))

	fs := http.FileServer(http.Dir("./ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fs))

	return mux
}

func (app *application) myhandler(w http.ResponseWriter, r *http.Request) http.Handler {
	return websocket.Handler(func(c *websocket.Conn) {

	})
}
