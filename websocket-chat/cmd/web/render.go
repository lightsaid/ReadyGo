package main

import (
	"fmt"
	"html/template"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, tplname string, data interface{}) {
	var err error
	ts := app.templateCache
	if !app.isProd {
		ts, err = app.genTemplate()
		if err != nil {
			fmt.Println("genTemplate error: ", err)
			w.WriteHeader(500)
			fmt.Fprintf(w, "创建模板失败")
			return
		}
	}

	t := ts.Lookup(tplname)
	if t == nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "模板不存在")
		return
	}

	t.Execute(w, data)
}

func (app *application) genTemplate() (*template.Template, error) {
	return template.ParseGlob("./ui/*.page.html")
}
