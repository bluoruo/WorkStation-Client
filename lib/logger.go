package lib

import (
	"fmt"
	"log"
	"os"
)

// 配置日志

var Logger *log.Logger
var ddnsLog *log.Logger
var repoLog *log.Logger
var statusLog *log.Logger
var putLog *log.Logger
var getLog *log.Logger

// var UpdateLog *log.Logger
var cmdLog *log.Logger

// 日志目录
var logPath = "./log"

// 日志文文件
var lgoFileBase = logPath + "/client.log"
var lgoFileRepo = logPath + "/repo.log"
var lgoFileDDns = logPath + "/ddns.log"
var lgoFileStatus = logPath + "/status.log"
var lgoFilePut = logPath + "/put.log"
var lgoFileGet = logPath + "/get.log"
var lgoFileCmd = logPath + "/cmd.log"

func init() {
	if !FileExists(logPath) {
		//常见日志文件夹
		err := os.Mkdir(logPath, 0777)
		if err != nil {
			fmt.Println("创建日志目录错误: ", err)
		}
	}

	//client.log
	logFile, err := os.OpenFile(lgoFileBase, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开日志文件异常: ", err)
	}
	//repo.log
	repoFile, err := os.OpenFile(lgoFileRepo, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开repo日志异常: ", err)
	}
	//ddns.log
	ddnsFile, err := os.OpenFile(lgoFileDDns, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开ddns日志异常: ", err)
	}
	//status.log
	statusFile, err := os.OpenFile(lgoFileStatus, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开status日志异常: ", err)
	}
	//put.log
	putFile, err := os.OpenFile(lgoFilePut, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开put日志异常: ", err)
	}
	//get.log
	getFile, err := os.OpenFile(lgoFileGet, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开get日志异常: ", err)
	}
	//cmd.log
	cmdFile, err := os.OpenFile(lgoFileCmd, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Panic("打开cmd日志异常: ", err)
	}

	Logger = log.New(logFile, "[Client]", log.Ldate|log.Ltime)
	repoLog = log.New(repoFile, "[REPO]", log.Ldate|log.Ltime)
	ddnsLog = log.New(ddnsFile, "[DDNS]", log.Ldate|log.Ltime)
	statusLog = log.New(statusFile, "[STATUS]", log.Ldate|log.Ltime)
	putLog = log.New(putFile, "[PUT]", log.Ldate|log.Ltime)
	getLog = log.New(getFile, "[GET]", log.Ldate|log.Ltime)
	cmdLog = log.New(cmdFile, "[CMD]", log.Ldate|log.Ltime)
}
