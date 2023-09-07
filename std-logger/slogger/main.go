package main

import (
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/exp/slog"
)

func main() {
	// firstLog()

	// secondLog()

	// thirdLog()

	// tourLog()

	mux := fifthLog()
	addr := "0.0.0.0:9000"
	slog.Info("server start on " + addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		slog.Error("server off", slog.String("err", err.Error()))
	}
}

// firstLog slog 初体验, 如何打印日志和输出额外参数
func firstLog() {
	slog.Debug("this is debug log")
	// 设置额外参数 slog.String(k,v)、slog.Int(k,v)、slog.Any(k,v) ... （这些就是slog.Attr）
	slog.Info("this is info log", slog.String("name", "lightsaid"))
	slog.Warn("this is warn log", slog.Any("aa", 100))
	slog.Error("this is error log")

	// slog.Level 默认有四个级别的日志
	// 其中slog 有一个 Default 默认实现的 slog.Logger, 并且默认级别时 Info，
	// 因此小于Info级别的不会输出，即Debug不会输出
	// slog 默认是文本输出，标准输出

	/*  slog.Attr 属性，k、v组成
	type Attr struct {
		Key   string
		Value Value
	}
	*/

	// 属性值分组输出
	slog.Error("属性值分组输出", slog.Group("book", slog.Int("id", 100), slog.String("name", "Go slog")))
}

// secondLog 创建一个logger实例，并设置slog默认值
func secondLog() {

	var out io.Writer = os.Stdout

	// NOTE: 如果向结合lumberjack输出到文件并实现日志分割，就是重新定义这个 &lumberjack.Logger{} 即可
	// 不多做解析，std-logger/logger 目录的代码里已实现
	// out = &lumberjack.Logger{}

	handler := slog.NewJSONHandler(
		out,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug, // 日志级别
			AddSource: true,            // 是否显示触发日志的位置
		},
	)
	l := slog.New(handler)

	// l.Debug("this is debug log")
	// l.Info("this is info log")
	// l.Warn("this is warn log")
	// l.Error("this is error log")

	// 设置为默认slog
	slog.SetDefault(l)
	firstLog()
}

// thirdLogger 使用 WithAttrs 添加默认字段
func thirdLog() {
	// 如果我希望每个日主输出都有默认字段，例如服务名字： server=orderServer
	// 这个该怎么做呢？
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug, // 日志级别
			AddSource: true,            // 是否显示触发日志的位置
		},
	).WithAttrs([]slog.Attr{
		slog.String("server", "OrderServer"),
		slog.Group("user", slog.Int("id", 100), slog.String("name", "lightsaid")),
	})
	l := slog.New(handler)

	l.Info("WithAttrs 使用")
}

// tourLog 使用 ReplaceAttr 替换属性，自此 slog.HandlerOptions 三个字段都清除明白
func tourLog() {
	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level:     slog.LevelDebug, // 日志级别
			AddSource: true,            // 是否显示触发日志的位置
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// 如果是 time 字段，则修改它的值
				if a.Key == "time" {
					a.Key = "date"
					a.Value = slog.StringValue(time.Now().Format("2006-01-02"))
				}
				return a
			},
		},
	)

	l := slog.New(handler)
	l.Info("ReplaceAttr的 使用")
}
