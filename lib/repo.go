package lib

import (
	"encoding/json"
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"os"
	"time"
)

/*
 * 基础信息相关
 */

type structBaseInfo struct {
	HostName string `json:"host_name"`
	HostId   string `json:"host_id"`
	OS       string `json:"os"`
	Arch     string `json:"arch"`
	Kernel   string `json:"kernel"`
	CPU      string `json:"cpu"`
	Mem      string `json:"mem"`
	DISK     string `json:"disk"`
}

var baseInfo = &structBaseInfo{}

var RepoServiceStatus = "none" //Repo服务状态

func init() {
	baseInfo.HostName, _ = os.Hostname() //HostName
	//OS 信息
	h, _ := host.Info()
	//fmt.Println("[Repo] debug all h host.Info: ", h)
	baseInfo.HostId = h.HostID // HostId
	if h.OS == "windows" {     // OS
		baseInfo.OS = h.Platform
	} else {
		baseInfo.OS = h.Platform + " " + h.PlatformVersion
	}
	baseInfo.Arch = h.KernelArch      //Arch
	baseInfo.Kernel = h.KernelVersion //Kernel
	//CPU 信息
	c, _ := cpu.Info()
	//fmt.Println("[Repo] debug all cpu.Info: ", c)
	baseInfo.CPU = c[0].ModelName
	//内存 信息
	m, _ := mem.VirtualMemory()
	mTotal := fmt.Sprintf("%.2f", float64(m.Total)/1024/1024/1024) //总大小 GB
	mUsedPer := fmt.Sprintf("%.2f", m.UsedPercent)                 //已使用 GB
	baseInfo.Mem = mTotal + "|" + mUsedPer
	//硬盘 信息
	d, _ := disk.Usage("/")
	dTotal := fmt.Sprintf("%.2f", float64(d.Total)/1024/1024/1024) //总大小 GB
	dUsed := fmt.Sprintf("%.2f", float64(d.Used)/1024/1024/1024)   //已使用 GB
	dUsedPer := fmt.Sprintf("%.2f", d.UsedPercent)                 //使用占比 %
	baseInfo.DISK = dTotal + "|" + dUsed + "|" + dUsedPer
}

// RunRepoService 运行 repo服务 [24小时更新一次]
func RunRepoService() {
	if RepoServiceStatus == "running" {
		repoLog.Println("RepoService has running.")
		return
	}
	for {
		if RepoServiceStatus == "stop" {
			RepoServiceStatus = "stopped"
			break
		}
		RepoServiceStatus = "running"
		SendWwsRepoInfo()
		time.Sleep(24 * time.Hour)
	}
}

// SendWwsRepoInfo 提交基本信息
func SendWwsRepoInfo() bool {
	body, err := HttpPost(serverConfig.Server+serverConfig.RepoApi,
		makeRepoInfoJson(),
		"json")
	if err != nil {
		repoLog.Println("Error 请求服务器端出错：", err)
	}
	data, st := CheckWssReturn(body)
	if st {
		repoLog.Println(" 提交info信息成功")
		return true
	} else {
		fmt.Println("API请求错误 Error:", data)
		repoLog.Println("Error API请求错误:", data)
	}
	return false
}

// SysInfo 基础信息
func SysInfo() {
	fmt.Println(makeRepoInfoJson())
}

// 构建 repo 接口 json信息
func makeRepoInfoJson() string {
	var repoInfo map[string]interface{}
	allInfo, _ := json.Marshal(baseInfo)
	_ = json.Unmarshal(allInfo, &repoInfo)
	delete(repoInfo, "host_id")
	repoInfo["ws_client_id"] = ClientBaseConfig.ID
	repoInfo["wsc_ver"] = AppVer
	reJson, _ := json.Marshal(repoInfo)
	return string(reJson)
}
