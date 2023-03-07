package main

import (
	userCom "client/command"
	"client/lib"
	"client/update"
	"fmt"
	"os"
)

const (
	confURL   = "https://ws.ofbao.com/v2/config.php"
	updateLog = "./log/update.log"
)

// 版本信息
var (
	BuildVersion string
	BuildCommit  string
	BuildTime    string
)

func init() {
	//单纯给批量编译提供一个版本号
	if len(os.Args) > 1 {
		if os.Args[1] == "version" {
			fmt.Printf("Version: %s\nBuild Time: %s\nGit Commit: %s\n",
				BuildVersion, BuildCommit, BuildTime)
			os.Exit(0)
		}
	}
	verInfo()                     //版权信息
	lib.AppVer = BuildVersion     //版本传递
	lib.ConfURL = confURL         //Config api 地址
	lib.LogFileUpdate = updateLog //update 日志文件
}

// 版权和版本信息
func verInfo() {
	fmt.Println("==============================================")
	fmt.Println("============= WorkStation Client =============")
	fmt.Println("======== ", BuildTime, "By Comanche Lab. ========")
	fmt.Println("================== Ver", BuildVersion, "==================")
	fmt.Println("============= Commit", BuildCommit, "=============")
	fmt.Println("==============================================")
}

// 初始化客户端
func initializationClient() {
	lib.Logger.Println("Start Client Ver", BuildVersion, "......")
	//初始化客户端
	lib.ServerConfig() //WorkStation Server 配置信息
	lib.StartClient()  //获取必要信息
	lib.Update()       //检查更新
}

func main() {
	update.DisUpdate = true //关闭启动更新
	//启动参数
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "pass": //鉴权启动
			fmt.Println("[启动方式] 鉴权启动")
			lib.Logger.Println("[启动方式] 鉴权启动")
			break
		case "-d": //后台启动
			lib.Logger.Println("[启动方式] 后台启动")
			break
		}
	}
	//默认执行
	go update.UpgradeFinish() //运行时检测upgrade是否存在(自更新必须！)
	initializationClient()    //初始化客户端

	//开启服务
	//等待服务器指令 1分钟执行一次
	lib.Logger.Println(" Start Wait Server Service...")
	go lib.RunWaitServerService()
	//后台启动 DDNS服务 10分钟执行一次
	lib.Logger.Println(" Start Client DDNS Service...")
	go lib.RunDDnsService()
	//后台启动 Status服务 30分钟执行一次
	lib.Logger.Println(" Start Client Status Service...")
	go lib.RunDDnsService()

	userCom.ClientShell() //执行命令

}
