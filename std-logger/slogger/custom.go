package main

import (
	"context"
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

// fifthLog 最后一个例子， 使用 Context 结合自定义slog.Handler 实现一个链路追踪
func fifthLog() http.Handler {
	handler := slog.NewJSONHandler(os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelDebug})

	myHandler := CustomHandler{
		Handler: handler,
	}

	l := slog.New(&myHandler)

	slog.SetDefault(l)

	mux := http.NewServeMux()
	mux.HandleFunc("/index", func(w http.ResponseWriter, r *http.Request) {
		data := make(map[string]string, 2)
		// 假设处理业务逻辑
		err := handlerIndex(r.Context())
		if err != nil {
			data["error"] = err.Error()
		} else {
			data["error"] = "ok"
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		// 获取 request_id
		requestID, _ := r.Context().Value("request_id").(string)
		data["request_id"] = requestID
		json.NewEncoder(w).Encode(data)
	})

	// 设置 request_id 中间件
	return setRequestID(mux)
}

func setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.Must(uuid.NewRandom())
		ctx := context.WithValue(r.Context(), "request_id", requestID.String())
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func handlerIndex(ctx context.Context) error {
	slog.InfoContext(ctx, "handler index request")
	time.Sleep(time.Second)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if r.Intn(10) > 5 {
		return errors.New("处理错误")
	}
	return nil
}

// CustomHandler 自定义 slog.Handler
type CustomHandler struct {
	slog.Handler
}

func (c *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if requestID, ok := ctx.Value("request_id").(string); ok {
		r.AddAttrs(slog.String("request_id", requestID))
	}

	return c.Handler.Handle(ctx, r)
}
