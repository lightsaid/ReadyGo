package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"readygo/wesocket-chat/db"
	"readygo/wesocket-chat/server"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

type application struct {
	port          int
	isProd        bool
	templateCache *template.Template
	serverWs      *server.ServerWs
}

func main() {
	var port int
	var isProd bool
	flag.IntVar(&port, "port", 3002, "web server port")
	flag.BoolVar(&isProd, "isProd", false, "is production env")
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%d", port)

	app := application{
		port:   port,
		isProd: isProd,
	}

	var err error
	app.templateCache, err = app.genTemplate()
	if err != nil {
		log.Fatal(err)
	}

	rdb = db.NewRedisClient()
	defer rdb.Close()

	// testRedis(cleint)

	srv := http.Server{
		Addr:           address,
		Handler:        app.routes(rdb),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		IdleTimeout:    time.Minute,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("server start on ", address)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// func testRedis(rdb *redis.Client) {
// 	ctx := context.Background()

// 	err := rdb.Set(ctx, "key", "value", 0).Err()
// 	if err != nil {
// 		panic(err)
// 	}

// 	val, err := rdb.Get(ctx, "key").Result()
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println("key", val)

// 	val2, err := rdb.Get(ctx, "key2").Result()
// 	if err == redis.Nil {
// 		fmt.Println("key2 does not exist")
// 	} else if err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println("key2", val2)
// 	}
// }
