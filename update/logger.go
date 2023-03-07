package update

import (
	"fmt"
	"log"
	"os"
)

var updateLog *log.Logger

func init() {
	var logPath = "./log"
	if !fileExists(logPath) {
		//常见日志文件夹
		err := os.Mkdir(logPath, 0777)
		if err != nil {
			fmt.Println("创建日志目录错误: ", err)
		}
	}
	//Update.log
	updateFile, err := os.OpenFile(logFileUpdate, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开update日志异常: ", err)
	}
	updateLog = log.New(updateFile, "[UPDATE]", log.Ldate|log.Ltime)
}
