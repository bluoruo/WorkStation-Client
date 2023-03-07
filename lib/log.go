package lib

import (
	"fmt"
	"os"
)

var LogFileUpdate string

func ReadLog(str string) {
	switch str {
	case "base":
		readLogLatest(lgoFileBase)
		break
	case "repo":
		readLogLatest(lgoFileRepo)
		break
	case "ddns":
		readLogLatest(lgoFileDDns)
		break
	case "status":
		readLogLatest(lgoFileStatus)
		break
	case "put":
		readLogLatest(lgoFilePut)
		break
	case "get":
		readLogLatest(lgoFileGet)
		break
	case "cmd":
		readLogLatest(lgoFileCmd)
		break
	case "update":
		readLogLatest(LogFileUpdate)
		break
	}
}

func readLogLatest(fileName string) {
	var maxLoad int64 = 1024
	// 打开文件
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(" [Log] 打开日志文件错误:", err)
	}
	defer file.Close()
	// 读入信息
	info, err := os.Stat(fileName)
	if err != nil {
		fmt.Println(" [Log] 读取日志文件错误:", err)
	}
	// 设置从哪里开始读取
	var buf []byte  //缓存区
	var start int64 //读取开始位置
	start = info.Size() - maxLoad
	if start < 0 { //不够 maxLoad
		start = 0
		buf = make([]byte, info.Size())
	} else {
		buf = make([]byte, maxLoad)
	}
	// 从缓存区中读取
	_, err = file.ReadAt(buf, start)
	if err == nil {
		fmt.Printf("%s", buf)
	}
}
