package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"readygo/simple-response-builder/respond"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_ = respond.NewResponse().
			Status(201).
			Header("Content-Type", "text/plain").
			String("Hello World!").
			Write(w)
	})

	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		book := Book{
			ID:     100,
			Title:  "Let's Go",
			Author: "Alex",
		}

		_ = respond.NewResponse().JSON(book).Status(http.StatusOK).Write(w)
	})

	// /file?filename=avatar2.png
	http.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		_ = respond.NewResponse().File(filename).Write(w)
	})

	// /base64?filename=avatar2.png
	http.HandleFunc("/base64", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Query().Get("filename")
		err := respond.NewResponse().Base64(filename).Write(w)
		if err != nil {
			log.Println(err)
		}
	})

	http.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		book := Book{
			ID:     100,
			Title:  "Let's Go",
			Author: "Alex",
		}
		t, _ := template.ParseFiles("./book.html")
		_ = respond.NewResponse().HTMLTemplate(t, book).Write(w)
	})

	addr := "0.0.0.0:9090"
	fmt.Println("server start on ", addr)
	http.ListenAndServe(addr, nil)
}
