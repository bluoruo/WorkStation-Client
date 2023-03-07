package lib

import (
	"client/update"
	"fmt"
	"strings"
	"time"
)

/*
 * 运行服务器发来的执行
 */

type structServerCom struct {
	Id   int
	Type int
	Com  string
}

var fromServerComQueue []structServerCom //队列
var runFromServerCom structServerCom     //运行中

var LockRunFromServerCom bool //运行服务器指令总开关
var lockConfig bool           //配置信息锁
var lockService bool          //服务锁
var lockCom bool              //指令锁
var lockTcp bool              //Tcp锁
var lockChannel bool          //通道锁

// RunWaitServerService 等待服务器指令
func RunWaitServerService() {
	for {
		if LockRunFromServerCom {
			strCom := fmt.Sprintf("ID: %v 类型: %v 内容: %s", runFromServerCom.Id, runFromServerCom.Type, runFromServerCom.Com)
			//检查运行指令的完整性
			if runFromServerCom.Com == "" {
				statusLog.Println("Error 指令信息不完整:", strCom)
				break
			}
			//运行
			switch runFromServerCom.Type {
			case 0: //更新配置相关
				if lockConfig {
					statusLog.Println("Wait 还未执行完毕")
					break
				}
				statusLog.Println("Run 配置指令")
				clientComFromServerConfig(runFromServerCom.Com)
				LockRunFromServerCom = false
				break
			case 1: //运行服务
				if lockService {
					statusLog.Println("Wait 还未执行完毕")
					break
				}
				statusLog.Println("Run 服务指令")
				clientComFromServerService(runFromServerCom.Com)
				LockRunFromServerCom = false
				break
			case 2: //运行指令
				if lockCom {
					statusLog.Println("Wait 还未执行完毕")
					break
				}
				statusLog.Println("Run 指令")
				clientComFromServerCom(runFromServerCom.Com)
				LockRunFromServerCom = false
				break
			case 3: //Tcp相关
				if lockTcp {
					statusLog.Println("Wait 还未执行完毕")
					break
				}
				statusLog.Println("Run Tcp相关指令")
				clientComFromServerTcp(runFromServerCom.Com)
				LockRunFromServerCom = false
				break
			case 5: //通道相关
				if lockChannel {
					statusLog.Println("Wait 还未执行完毕")
					break
				}
				statusLog.Println("Run 通道相关指令")
				clientComFromServerChannel(runFromServerCom.Com)
				LockRunFromServerCom = false
				break
			default:
				statusLog.Println("Error 无效的指令类型:", strCom)
				LockRunFromServerCom = false
			}
		}
		//statusLog.Println("Debug 没有指令,等待10秒.")
		time.Sleep(10 * time.Second)
	}
	statusLog.Println("Exit 退出指令服务！")
}

// 配置
func clientComFromServerConfig(strCom string) {
	lockConfig = true //加锁
	switch strCom {
	case "config": //更新接口配置文件信息
		statusLog.Println(" 更新ServerConfig信息")
		updateServerConfig()
		break
	case "ddns": //更新ddns接口配置信息
		UpdateWwsDDnsConfig()
		statusLog.Println(" 更新ddns接口信息")
		break
	case "update": //更新 (特殊！堵塞执行，然后然后执行更新)
		statusLog.Println(" 更新程序")
		update.DisUpdate = false
		update.RunClientUpdateCom()
		break
	}
	lockConfig = false //解锁
}

// 服务
func clientComFromServerService(strCom string) {
	arrCom := strings.Split(strCom, ";")
	lockService = true //加锁
	switch arrCom[0] {
	case "repo": //repo服务
		switch arrCom[1] {
		case "start": //启动服务
			statusLog.Println(" 启动REPO服务...")
			go RunRepoService()
			break
		case "stop": //停止服务
			statusLog.Println(" 停止REPO服务.")
			RepoServiceStatus = "stop"
			break
		case "once": //执行一次
			statusLog.Println(" 执行一次REPO服务.")
			break
		}
		break
	case "ddns": //ddns服务
		switch arrCom[1] {
		case "start":
			statusLog.Println(" 启动DDNS服务...")
			go RunDDnsService()
			break
		case "stop":
			statusLog.Println(" 停止DDNS服务.")
			DDnsServiceStatus = "stop"
			break
		case "once":
			statusLog.Println(" 执行一次DDNS服务.")
			break
		}
		break
	case "status": //status服务 (此服务不能停止！)
		switch arrCom[1] {
		case "start": //启动服务
			statusLog.Println(" 启动Status服务...")
			go RunDDnsService()
			break
		case "stop": //停止服务
			//DDnsServiceStatus = "stop"
			statusLog.Println(" error 当前版本不允许停止Status服务")
			break
		case "once": //执行一次
			statusLog.Println(" 执行一次STATUS服务.")
			break
		}
		break
	}
	lockService = false //解锁
}

// 指令  [com;[sys]ls -l] [com;[sys]dir C:\] [com;[user]md5sum.exe -f notepad.exe]
func clientComFromServerCom(strCom string) {
	arrCom := strings.Split(strCom, ";")
	lockCom = true //加锁
	switch arrCom[0] {
	case "com":
		statusLog.Println(" 执行：", arrCom[1])
		go RunCommCom(arrCom[1])
		break
	case "put": //上传文件
		statusLog.Println(" 上传文件：", arrCom[1])
		break
	case "get": //下载文件
		statusLog.Println(" 下载文件：", arrCom[1])
		break
	}
	lockCom = false //解锁
}

// Tcp
func clientComFromServerTcp(strCom string) {
	arrCom := strings.Split(strCom, ";")
	lockTcp = true //加锁
	switch arrCom[0] {
	case "tcp": //tcp
		statusLog.Println(" TCP通讯：", arrCom[1])
		break
	case "udp": //udp
		statusLog.Println(" UDP通讯：", arrCom[1])
		break
	}
	lockTcp = false //解锁
}

// 通道
func clientComFromServerChannel(strCom string) {
	arrCom := strings.Split(strCom, ";")
	lockChannel = true //加锁
	switch arrCom[0] {
	case "tcp": //tcp方式
		statusLog.Println(" TCP通道：", arrCom[1])
		break
	case "udp": //udp 方式
		statusLog.Println(" UDP通道：", arrCom[1])
		break
	}
	lockChannel = false //解锁
}
