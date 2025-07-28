package utils

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func init() {
	logFile, err := os.OpenFile(RuntimeLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	// 创建一个替换时间的函数
	customTimeFormat := "2006-01-02 15:04:05"
	replaceTime := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(customTimeFormat))
		}
		return a
	}

	Logger = slog.New(slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource:   true,           // 记录错误位置
		Level:       slog.LevelInfo, // 设置日志级别
		ReplaceAttr: replaceTime,
	}))
}

func InitAccessLogger() *os.File {
	f, err := os.OpenFile(AccessLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Logger.Error("创建请求日志失败")
		panic("创建请求日志失败")
	}
	return f
}
