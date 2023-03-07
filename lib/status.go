package lib

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"strings"
	"time"
)

/*
 * 状态信息相关
 */

// status接口 提交数据结构
type structStatusInfo struct {
	WsClientId int    `json:"ws_client_id"`
	Cpu        string `json:"cpu"`
	CpuTop     string `json:"cpu_top"`
	Mem        string `json:"mem"`
	MemTop     string `json:"mem_top"`
	Hd         string `json:"hd"`
	Net        string `json:"net"`
	Status     int    `json:"status"`
}

// status接口 返回数据结构
type structStatusReturnJson struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Id   int    `json:"id"`
		Type int    `json:"type"`
		Com  string `json:"com"`
	} `json:"data"`
}

var statusInfo = &structStatusInfo{}
var statusReturnJson = &structStatusReturnJson{}
var StatusServiceStatus = "none" //Status服务状态

// 构建 status 接口 json数据
func makeStatusInfoJson() string {
	//CPU 信息
	cpuPercents, _ := cpu.Percent(time.Second, false)
	var totalCpuPre float64
	var i int
	for i = 0; i < len(cpuPercents); i++ {
		totalCpuPre = totalCpuPre + cpuPercents[i]
	}
	//fmt.Println("[Status] debug CPU 负载:", cpuPercent)
	cpuPer := fmt.Sprintf("%.4f", totalCpuPre/float64(i))

	//内存 信息
	m, _ := mem.VirtualMemory()
	memoryTotal := fmt.Sprintf("%.2f", float64(m.Total)/1024/1024/1024) //总大小 GB
	memoryUsed := fmt.Sprintf("%.2f", float64(m.Used)/1024/1024/1024)   //已使用大小 GB
	memoryUsedPer := fmt.Sprintf("%.2f", m.UsedPercent)                 //已使用百分比 %
	//fmt.Println("[Status] debug 内存信息: Total->", memoryTotal, "G Used->", memoryUsed, " G Used Per->", memoryUsedPer)
	memory := formatInfoString(memoryTotal) + "|" + formatInfoString(memoryUsed) + "|" + formatInfoString(memoryUsedPer)
	//fmt.Println("[Status] debug 内存信息:", memory)

	//硬盘 信息
	d, _ := disk.Usage("/")
	hdTotal := fmt.Sprintf("%.2f", float64(d.Total)/1024/1024/1024) //总大小 GB
	hdUsed := fmt.Sprintf("%.2f", float64(d.Used)/1024/1024/1024)   //已使用 GB
	hdUsedPer := fmt.Sprintf("%.2f", d.UsedPercent)                 //使用占比 %
	//fmt.Println("HD: Total->", diskTotal, "G Used->", disUsed, "G Used Per->", disUsedPer)
	hd := formatInfoString(hdTotal) + "|" + formatInfoString(hdUsed) + "|" + formatInfoString(hdUsedPer)
	//fmt.Println("[Status] Debug 内存信息:", hd)

	//网络 信息

	statusInfo.WsClientId = ClientBaseConfig.ID
	statusInfo.Cpu = cpuPer
	statusInfo.CpuTop = "cpu负载"
	statusInfo.Mem = memory
	statusInfo.MemTop = "内存相关"
	statusInfo.Hd = hd
	statusInfo.Net = "none"
	statusInfo.Status = 1
	reJson, _ := json.Marshal(statusInfo)
	return string(reJson)
}

// RunStatus 后台运行 status 服务
func RunStatus() {
	if DDnsServiceStatus == "running" {
		fmt.Println("Status Service has running.")
		return
	}
	statusLog.Println("Start Status Service....")
	for {
		if StatusServiceStatus == "stop" { //停止运行
			StatusServiceStatus = "stopped"
			break
		}
		//加载最新的配置信息
		iniLoad()
		//开始运行
		StatusServiceStatus = "running"
		if !SendWwsStatusInfo() {
			statusLog.Println("向服务器提交信息出错了")
		}
		//30分钟提交一次系统状态
		time.Sleep(30 * time.Minute)
	}
}

// SendWwsStatusInfo 提交Status信息
func SendWwsStatusInfo() bool {
	body, err := HttpPost(serverConfig.Server+serverConfig.StatusApi,
		makeStatusInfoJson(),
		"json")
	if err != nil {
		fmt.Println("[Status] 请求服务器端出错：", err)
		statusLog.Println("请求服务器端出错", err)
	}
	err = json.Unmarshal(body, statusReturnJson)
	if err != nil {
		fmt.Println("[Status] 返回信息错误 Error:", err)
		statusLog.Println("返回信息错误 Error:", err)
	}

	if statusReturnJson.Code == "0" {
		if statusReturnJson.Msg == "ok" {
			return true
		}
		if statusReturnJson.Msg == "success" {
			//fmt.Println(statusReturnJson.Data)
			statusLog.Println("Debug 服务器要求执行任务！")
			// 去执行指令
			if LockRunFromServerCom { //已有指令在运行则加入队列
				statusLog.Println("Debug 服务器指令锁定,压栈等待！")
				oneCom := structServerCom{
					Id:   statusReturnJson.Data.Id,
					Type: statusReturnJson.Data.Type,
					Com:  statusReturnJson.Data.Com}
				// 如果指令已经运行 则 加入排队
				fromServerComQueue = append(fromServerComQueue, oneCom)
			} else { //直接执行
				statusLog.Println("Debug 10秒后执行服务器指令！")
				LockRunFromServerCom = true
				runFromServerCom = structServerCom{
					Id:   statusReturnJson.Data.Id,
					Type: statusReturnJson.Data.Type,
					Com:  statusReturnJson.Data.Com}
			}
		}
	} else {
		fmt.Println("[Status] API错误:", statusReturnJson.Msg)
	}
	return false
}

// StatusInfo 显示Status信息
func StatusInfo() {
	fmt.Println(makeStatusInfoJson())
}

// 去掉最后一位的点
func formatInfoString(str string) string {
	//去掉末尾的0
	str = strings.TrimRight(str, "0")
	//去掉最后一位的。
	if str[len(str)-1:] == "." {
		str = str[:len(str)-1]
	}
	return str
}
