package util

import (
	"encoding/base64"
	"errors"
	"net/http"
)

var (
	ErrValueTooLong = errors.New("cookie值太长了")
	ErrInvalidValue = errors.New("无效的cookie值")
)

// 每个Cookie最大值4k
const cookieMaxSize = 4096

// Write 简单地写入Cookie，仅做URL编码
func WriteCookie(w http.ResponseWriter, cookie http.Cookie) error {
	// 将 cookie 值编码
	cookie.Value = base64.URLEncoding.EncodeToString([]byte(cookie.Value))

	if len(cookie.String()) > cookieMaxSize {
		return ErrValueTooLong
	}

	http.SetCookie(w, &cookie)

	return nil
}

// Read 简单地读取cookie，仅做URL解码
func ReadCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	value, err := base64.URLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return "", ErrInvalidValue
	}

	return string(value), nil
}
