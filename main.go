package main

import (
	"gmeroblog/model"
	"gmeroblog/routes"
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

// 初始化
func init() {
	// 记录核心日志
	logFile := &lumberjack.Logger{
		Filename:   "log/core.log",
		MaxSize:    1,     // 文件大小MB
		MaxBackups: 5,     // 最大保留日志文件数量
		MaxAge:     28,    // 保留天数
		Compress:   false, // 是否压缩
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Llongfile | log.Ltime | log.Ldate)
}

func main() {
	model.InitDb()
	routes.InitRouter()
}
