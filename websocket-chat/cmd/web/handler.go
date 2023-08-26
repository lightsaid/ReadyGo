package main

import (
	"net/http"
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
