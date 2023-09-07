package main

import (
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type F map[string]interface{}

func main() {
	// firstLog()
	// secondLog()
	// thirdLogger()
	// fourLogger()

	// 使用自定义JSONLogger
	slog := NewJSONLogger(os.Stderr, InfoLevel)
	slog.Info("this is info log", nil)
	slog.Debug("this is debug log", nil) // 不会输出
	slog.Error("this error log ", F{"id": 100, "name": "诸葛亮"})
	slog.Error("this error log ", F{"id": 200, "name": "曹操"}, true)

	// 结合 lumberjack 日志分割，不多举例，参考 fourLogger
	// slog2 := NewJSONLogger(&lumberjack.Logger{}, InfoLevel)

}

// firstLog 简单使用，适合于项目
func firstLog() {
	// 使用 log.New 创建一个新的log实例，可以设置日志默认前缀和一些flag（如日期、时间、文件）
	logger := log.New(os.Stdout, "INFO ", log.Ldate|log.Ltime|log.Llongfile)

	// INFO 2023/09/06 23:42:00 /home/lightsaid/go/src/ReadyGo/std-logger/logger/main.go:12: this is one log
	logger.Println("this is one log")
}

// secondLog 在firstLog的基础上将日志输入到文件
func secondLog() {

	file, err := os.OpenFile("./logger/log.out", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	logger := log.New(file, "INFO ", log.Ldate|log.Ltime|log.Llongfile)

	logger.Println("this is two log")
}

// thirdLogger 在 secondLog 基础上监听日志文件大小，实现轮循日志分割
func thirdLogger() {
	file, err := os.OpenFile("./logger/log2.out", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	logger := log.New(file, "INFO ", log.Ldate|log.Ltime|log.Llongfile)

	var mu sync.RWMutex

	go func() {
		for {
			finfo, err := file.Stat()
			if err != nil {
				log.Println(err)
				break
			}
			time.Sleep(1 * time.Second)

			// 假设日志文件大于100kb即分割日志
			if finfo.Size() > 100 {
				mu.Lock()

				ts := time.Now().Format("2006-01-02 15:04:05")
				name := "log_" + ts + ".log"
				newFile, err := os.OpenFile(name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
				if err != nil {
					log.Println(err)
					break
				}
				// 重新设置输出目标
				logger.SetOutput(newFile)

				mu.Unlock()
			}
		}
	}()

	var i = 0
	for {
		i++
		logger.Println("this is two log")
		if i > 10000 {
			break
		}
		time.Sleep(100 * time.Microsecond)
	}
}

// fourLogger 在 secondLog 的基础上 + lumberjack 实现日志分割
func fourLogger() {

	logWriter := &lumberjack.Logger{
		Filename:   "./logger/foo.log",
		MaxSize:    1, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	}

	logger := log.New(logWriter, "INFO ", log.Ldate|log.Ltime|log.Llongfile)

	ticker := time.NewTicker(3 * time.Second)
	f := true
	go func() {
		<-ticker.C
		f = false
	}()
	for {
		if !f {
			break
		}
		logger.Println("this is foo log")
		time.Sleep(50 * time.Microsecond)
	}
}
