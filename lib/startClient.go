package lib

import (
	"bufio"
	"fmt"
	"os"
)

/*
 * 初始化客户端
 */

// Client端 name 检查
func checkLocalName() bool {
	name := ClientBaseConfig.Name
	if name == "" { //不存在设置名
		fmt.Println("请设置一个设备识别名:")
		name, _ = bufio.NewReader(os.Stdin).ReadString('\n')
		name = Strip(name, "/ \r\n") //去掉 \r\n
		if name == "" || checkBaseName(name) == false {
			return false
		}
	}
	return true
}

// StartClient 入口
func StartClient() {
	fmt.Println("Start WorkStation Client base......")
	ServerConfig()        //WorkStation Server 配置信息
	if checkLocalName() { //识别名检测
		fmt.Println(" [初始化] 识别名 OK!")
		Logger.Println(" [初始化] 识别名 OK!")
	} else {
		fmt.Println(" [初始化] Error 设置识别名错误！")
		Logger.Println(" [初始化] Error 设置识别名错误！")
		return
	}
	if GetWwsBaseConfig() { //WorkStation Server Base信息
		fmt.Println(" [初始化] 获取基本信息 OK!")
		Logger.Println(" [初始化] 获取基本信息 OK!")
	} else {
		fmt.Println(" [初始化] Error 获取基本错误!")
		Logger.Println(" [初始化] Error 获取基本信息完毕！")
		return
	}
	if UpdateWwsDDnsConfig() { //WorkStation Server DDns接口信息
		fmt.Println(" [初始化] 获取DDNS接口信息 OK!")
		Logger.Println(" [初始化] 获取DDNS接口信息 OK!")
	} else {
		Logger.Println(" [初始化] Error 获取DDNS接口信息错误！")
		return
	}
	if SendWwsRepoInfo() { //WorkStation Server 提交本机信息
		fmt.Println(" [初始化] 本机信息检查 OK!")
		Logger.Println(" [初始化] 本机信息检查 OK!")
	} else {
		Logger.Println(" [初始化] Error 本机信息信息错误！")
		return
	}
	fmt.Println("Start WorkStation Client Finished.")
	Logger.Println("Start WorkStation Client Finished.")
	return
}
