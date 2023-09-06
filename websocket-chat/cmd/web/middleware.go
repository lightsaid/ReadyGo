package main

import (
	"fmt"
	"log"
	"net/http"
	"readygo/wesocket-chat/db"
	"readygo/wesocket-chat/model"
	"time"
)

type middlewareHandler func(w http.ResponseWriter, r *http.Request, u *model.User)

// middlewareAuth 从cookie获取sessionID 查询用户信息，提供给下一个handler处理函数
func (app *application) middlewareAuth(next middlewareHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session")
		if err != nil {
			if err != http.ErrNoCookie {
				log.Println("get session cookie error: ", err)
			}
			next(w, r, nil)
			return
		}

		if cookie.Value == "" {
			next(w, r, nil)
			return
		}

		ss, err := db.RepoObj.SessionRepo.Get(cookie.Value)
		if err != nil {
			log.Println("from SessionRepo get session error: ", err)
			next(w, r, nil)
			return
		}

		// cookie 过期
		if ss.Expires.Before(time.Now()) {
			// 清理cookie
			http.SetCookie(w, &http.Cookie{Name: "session", MaxAge: -1})
			next(w, r, nil)
			return
		}

		user, err := db.RepoObj.UserRepo.GetByID(ss.UserID.String())
		if err != nil {
			log.Println("from UserRepo.GetByID get user error: ", err)
			next(w, r, nil)
			return
		}

		fmt.Println("auth -> user: ", user.Nickname, ss.Expires)
		next(w, r, user)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if f := recover(); f != nil {
				log.Printf("Panic: %v\n", f)
			}
		}()

		next.ServeHTTP(w, r)

	})
}
