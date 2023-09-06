package respond

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

type responseType int

const (
	textResponse responseType = iota
	jsonResponse
	fileResponse
	htmlResponse
	base64Response
)

type response struct {
	status       int
	headers      map[string]string
	data         interface{}
	htmlTemplate *template.Template
	responseType
}

func NewResponse() *response {
	return &response{
		status:  200,
		headers: make(map[string]string),
	}
}

func (r *response) Status(statusCode int) *response {
	r.status = statusCode
	return r
}

func (r *response) Header(key, val string) *response {
	r.headers[key] = val
	return r
}

func (r *response) String(message string) *response {
	r.responseType = textResponse
	r.Header("Content-Type", "text/plain")
	r.data = message
	return r
}

func (r *response) JSON(data interface{}) *response {
	r.responseType = jsonResponse
	r.Header("Content-Type", "application/json")
	r.data = data
	return r
}

func (r *response) File(filepath string) *response {
	r.responseType = fileResponse
	r.data = filepath
	return r
}

func (r *response) HTMLTemplate(htmlTpl *template.Template, data interface{}) *response {
	r.responseType = htmlResponse
	r.Header("Content-Type", "text/html")
	r.htmlTemplate = htmlTpl
	r.data = data
	return r
}

func (r *response) Base64(filename string) *response {
	r.responseType = base64Response
	r.data = filename
	return r
}

func (r *response) Write(w http.ResponseWriter) error {
	var err error

	// 添加 Header
	for k, v := range r.headers {
		fmt.Println(k, " = ", v)
		w.Header().Add(k, v)
	}

	// // 最后才能写入状态，一旦写入状态，再添加Header头也是无效的
	// w.WriteHeader(r.status)

	switch r.responseType {
	case textResponse:
		w.WriteHeader(r.status)
		_, err = w.Write([]byte(r.data.(string)))
	case jsonResponse:
		w.WriteHeader(r.status)
		err = json.NewEncoder(w).Encode(r.data)
	case fileResponse:
		var bs []byte
		bs, err = os.ReadFile(r.data.(string))
		if err != nil {
			return err
		}
		w.WriteHeader(r.status)
		_, err = w.Write(bs)
	case htmlResponse:
		w.WriteHeader(r.status)
		err = r.htmlTemplate.Execute(w, r.data)
	case base64Response:
		var bs []byte
		bs, err = os.ReadFile(r.data.(string))
		if err != nil {
			return err
		}

		var result string
		mimeType := http.DetectContentType(bs)

		switch mimeType {
		case "image/jpeg":
			result += "data:image/jpeg;base64,"
		case "image/png":
			result += "data:image/png;base64,"
			// TODO: 其他类型
			// ...
		}
		// w.Header().Add("Content-Type", mimeType)
		result += base64.StdEncoding.EncodeToString(bs)
		w.WriteHeader(r.status)
		w.Write([]byte(result))
	}

	return err
}
