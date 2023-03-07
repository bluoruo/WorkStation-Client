package command

import (
	"bufio"
	"client/lib"
	"client/update"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ClientShell 入口
func ClientShell() {
	whileCmd()
}

/*
 * 输出帮助信息
 */
func helpInfo() {
	fmt.Println("Option:")
	fmt.Println("  cat <ddns/repo/status>      view service info")
	fmt.Println("  set <ddns>                  setup service")
	fmt.Println("  run <ddns/repo/status>      run service")
	fmt.Println("  log <base/repo/ddns/status> view logs")
	fmt.Println("  put <to> <file>             upload file (远程目录只需要指定文件夹)")
	fmt.Println("  get <remote file>           download file (下载到程序所在文件夹)")
	fmt.Println("  update                      check and upgrade latest version")
	fmt.Println("  quit                        exit client")
	fmt.Println("  cmd <shell> exec command")
}

/*
 * 读取用户输入命令
 * @return string
 */
func readIOLine() string {
	line, err := bufio.NewReader(os.Stdin).ReadString(byte('\n'))
	if err != nil && err != io.EOF {
		//fmt.Println("[Cmd] debug 输入错误, Error:", err)
	}
	/*
	 * 问题描述：
	 *   在不判断是否是后台运行时，Linux后台会不停请求,
	 *   会造成写入大量日志或者IO占用。
	 * 解决方法:
	 *   1.用换行符来区分是否是用户在输入,非用户输入等待1小时。
	 *   2.后台运行时，使用TCP或者HTTP服务代替命令行输入。
	 * Ver 0.26	2023/03/05 by Harry.Ren
	 */
	//判断是否是后台运行
	if strings.Contains(line, "\n") {
		//用户运行直接返回
		return strings.TrimSpace(line)
	} else { //后台运行的话 堵塞1小时
		fmt.Println("[Cmd] 后台运行，堵塞1小时")
		time.Sleep(1 * time.Hour)
	}
	return ""
}

/*
 * 循环读取用户输入 (基础命令)
 */
func whileCmd() {
	for {
		fmt.Print("wsc-> ")  //命令提示符
		text := readIOLine() //接受用户输入命令
		if text == "" {      //空指令处理
			continue
		}
		//处理用户命令
		switch text {
		case "help": //帮助
			helpInfo()
			break
		case "update": //更新
			runUpdate()
			break
		case "ver": //查看版本
			fmt.Println("软件版本:", lib.AppVer, " By Comanche Lab.")
			break
		case "quit": //退出
			return
		default:
			if runClientCmd(text) { //客户端服务指令
				continue
			}
		}
	}
}

/*
 * 执行用户输入命令 (带有参数的命令)
 * @param text string
 * @return bool
 */
func runClientCmd(text string) bool {
	//带有参数的命令
	if strings.Contains(text, " ") {
		arrText := strings.Fields(text)
		if len(arrText) > 0 {
			switch arrText[0] {
			case "put": //上传文件
				runClientPUT(arrText)
				break
			case "get": //下载文件
				runClientGET(arrText)
				break
			case "log": //日志
				runClientLOG(arrText)
				break
			case "cmd": //系统命令
				runClientCMD(arrText)
				break
			case "cat", "set", "run", "stop": //Client服务名命令
				runClientService(arrText)
				break
			}
		}
	}
	//return true 继续接收命令输入
	return true
}

/*
 * 执行 Client 命令
 */

// Client 发送文件
func runClientPUT(arrText []string) {
	var putPort = ":28169" //发送文件使用的端口
	if arrText[1] == "to" && len(arrText) == 4 {
		//发送文件客户端
		lib.PutClient(arrText[2]+putPort, arrText[3])
	} else if arrText[1] == "listen" {
		//发送文件服务端
		lib.PutServer("0.0.0.0" + putPort)
	} else {
		fmt.Println("put <option>:")
		fmt.Println("put <to> <host> <file>")
		fmt.Println("put <listen>")
	}
}

// Client 下载文件
func runClientGET(arrText []string) bool {
	return true
}

// Client 查看日志
func runClientLOG(arrText []string) bool {
	lib.ReadLog(arrText[1])
	return true
}

// Client 系统命令
func runClientCMD(arrText []string) bool {
	switch arrText[1] {
	case "shell":
		if arrText[2] != "" {
			fmt.Println(arrText[1], "to", arrText[2])
			lib.TestCmd(arrText[2])
		} else {
			fmt.Println("cmd shell <host:port>  获取一个反弹Shell")
		}
		break
	case "channel":
		if arrText[2] != "" && arrText[3] != "" && arrText[4] != "" {
			fmt.Println(arrText[3], "to", arrText[4])
		} else {
			fmt.Println("cmd channel <option> <local host:port> <remote host:port>  获取一个反向通道")
		}
		break
	case "socks":
		if arrText[2] != "" && arrText[3] != "" {
			fmt.Println("Socks5 to", arrText[4])
		} else {
			fmt.Println("cmd channel <option> <remote host:port>  获取一个反向Socks5代理")
		}
	}

	return true
}

// Client 连接服务器
func connServer(arrText []string) {
	var serverPort = ":28169"
	lib.ConnectTCPServer(arrText[1] + serverPort)
}

// Client 服务器端
func startServer(arrText []string) {
	var serverPort = ":28169"
	lib.ListenTCPServer("0.0.0.0" + serverPort)
}

// Client 软件更新
func runUpdate() {
	update.DisUpdate = false
	lib.Update()
}

/*
 * 执行 Client 服务 命令
 */
func runClientService(arrText []string) {
	switch arrText[1] {
	case "config":
		runServiceCONFIG(arrText[0])
		break
	case "ddns":
		runServiceDDNS(arrText[0])
		break
	case "repo":
		runServiceREPO(arrText[0])
		break
	case "status":
		runServiceStatus(arrText[0])
		break
	}
}

// Service 执行CONFIG相关
func runServiceCONFIG(text string) {
	switch text {
	case "cat":
		//显示当前配置信息
		break
	case "set":
		//修改参数
		break
	case "run":
		//重新获取config信息
		break
	default:
		println("无效参数")
	}
}

// Service 执行DDNS相关
func runServiceDDNS(text string) {
	switch text {
	case "cat":
		//本地IP信息
		go lib.DDnsInfo()
		break
	case "set":
		if lib.GetWwsBaseConfig() { //获取Base配置
			if lib.UpdateWwsDDnsConfig() { //获取DDNS接口配置
				fmt.Println("获取到DDNS新的接口配置.")
				fmt.Println("使用命令 run ddns 启动DDNS服务")
				break
			}
		}
		fmt.Println("服务器还未更新配置，请通知管理员后重试！")
		break
	case "run": //后台运行
		go lib.RunDDnsService()
		break
	case "stop": //停止运行
		lib.DDnsServiceStatus = "stop"
		break
	default:
		fmt.Println("无效参数")
	}
}

// Service 执行REPO相关
func runServiceREPO(text string) {
	switch text {
	case "cat":
		lib.SysInfo() //显示基础信息
		break
	case "run":
		go lib.SendWwsRepoInfo() //后台 提交info
		break
	default:
		fmt.Println("无效参数")
	}
}

// Service 执行Status相关
func runServiceStatus(text string) {
	switch text {
	case "cat":
		lib.StatusInfo() //显示status信息
		break
	case "set":
		go lib.SendWwsStatusInfo() //单次提交
		break
	case "run":
		go lib.RunStatus() //后台运行
		break
	case "stop": //暂不支持停止status服务
		break
	default:
		fmt.Println("无效参数")
	}
}
