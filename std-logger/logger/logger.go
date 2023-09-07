package main

import (
	"encoding/json"
	"io"
	"runtime/debug"
	"sync"
	"time"
)

// 在 main.go 中编写4个方法介绍标准库 log 使用。

// 在此文件中将实现一个带有等级的JSON格式日志

// NOTE: 思路 伪代码，
// var out io.Writer
// 最终输出日志方法
// out.Write([]byte("日志"))

// 这里将日志设置为4个级别分别是：debug、info、error、panic，分别实现每个级别的日志输出
// 四个等级的日志分别收集基础信息，最终统一方法打印
// 定义一个结构体，保存日志的最小级别，仅当日志大于或等于最小级别才输出
// 定义结构体设置日志的一些基础信息，如message、time、并可以扩充 map[string]interface{}
// 最后是否要打印日志堆栈信息，可以使用 debug.Stack() 获取

type Level uint8

// 设置为4个等级的日志
const (
	DebugLevel Level = iota // 0
	InfoLevel               // 1
	ErrorLevel              // 2
	PanicLevel              // 3
)

func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	default:
		return ""
	}
}

// 定义一个JSON logger 结构体，实现JSON日志输出
type JSONLogger struct {
	out       io.Writer
	mineLevel Level
	timeFomat string
	mutex     sync.Mutex
}

func NewJSONLogger(out io.Writer, mineLevel Level) *JSONLogger {
	return &JSONLogger{
		out:       out,
		mineLevel: mineLevel,
		timeFomat: time.RFC3339,
	}
}

// SetTimeFormat 设置日期输出格式
func (l *JSONLogger) SetTimeFormat(tf string) {
	l.timeFomat = tf
}

func (l *JSONLogger) Debug(msg string, fields map[string]interface{}, trace ...bool) {
	l.print(DebugLevel, msg, fields, trace...)
}

func (l *JSONLogger) Info(msg string, fields map[string]interface{}, trace ...bool) {
	l.print(InfoLevel, msg, fields, trace...)
}

func (l *JSONLogger) Error(msg string, fields map[string]interface{}, trace ...bool) {
	l.print(ErrorLevel, msg, fields, trace...)
}

func (l *JSONLogger) Panic(msg string, fields map[string]interface{}, trace ...bool) {
	l.print(PanicLevel, msg, fields, trace...)
}

// print 统一打印日志函数
func (l *JSONLogger) print(level Level, message string, fields map[string]interface{}, traces ...bool) (int, error) {
	if level < l.mineLevel {
		return 0, nil
	}

	var trace bool
	if len(traces) > 0 {
		trace = traces[0]
	}

	var slog = struct {
		Level   string                 `json:"level"`
		Time    string                 `json:"time"`
		Message string                 `json:"message"`
		Fields  map[string]interface{} `json:"fields,omitempty"` // 其他字段
		Trace   string                 `json:"trace,omitempty"`  // 堆栈信息
	}{
		Level:   l.mineLevel.String(),
		Time:    time.Now().Format(l.timeFomat),
		Message: message,
		Fields:  fields,
	}

	if trace {
		slog.Trace = string(debug.Stack())
	}

	var line []byte
	line, err := json.Marshal(slog)
	if err != nil {
		line = []byte(ErrorLevel.String() + ": json 序列化日志失败： " + err.Error())
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()

	// 添加换行符号
	line = append(line, '\n')

	return l.out.Write([]byte(line))
}
