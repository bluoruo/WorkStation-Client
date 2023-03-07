package update

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

/*
 * upgrade 相关
 */

// 检查ws_upgrade.exe文件
func checkUpgradeFile() bool {
	var path = "./"        //目录
	var name = upgradeName //去掉.exe的程序名
	if sysOS == "windows" {
		name = strings.TrimSuffix(upgradeName, ".exe")
	}
	if fileExists(path + upgradeName) {
		return true
	} else { //不存在，开始下载
		fmt.Println("[Upgrade] 开始下载更新程序" + name + ".zip...")
		updateLog.Println("开始下载更新程序" + name + ".zip...")
		err := downUpgradeZip(path, name)
		if err != nil {
			//fmt.Println("[Upgrade] Error 下载"+name+".zip失败:", err)
			updateLog.Println("Error 下载"+name+".zip失败:", err)
			return false
		} else {
			fmt.Println("[Upgrade] " + name + ".zip下载完成,解压zip文件...")
			updateLog.Println(name, ".zip下载完成,解压zip文件...")
			err = safeUnZip(path, name)
			if err != nil { //解压错误
				//fmt.Println("[Upgrade] debug 解压文件出错！", err)
				return false
			}
			fmt.Println("[Upgrade] " + name + ".zip解压完成！")
			updateLog.Println(name + ".zip解压完成！")
		}
	}
	return true
}

// 下载upgrade.zip
func downUpgradeZip(path, name string) error {
	err := downloadFile("https://ws.ofbao.com/down/base/"+name+".zip", path+name+".zip")
	return err
}

// Upgrade 是否运行
func checkUpgradeRunStatus(cStr string) bool {
	var str string
	var err error
	var st = false
	for i := 0; i < 20; i++ {
		//fmt.Println("[Upgrade] debug 002-1 check upgrade run.")
		str, err = startTcpClient("by Client")
		if err == nil {
			//fmt.Println("[Upgrade] debug 001 return Msg:", str)
			if str == cStr {
				//fmt.Println("[Upgrade] debug 002-2 has OK!")
				st = true
				break
			} else {
				//fmt.Println("[Upgrade] 端口开启,不是" + upgradeName)
				updateLog.Println("端口开启,不是" + upgradeName)
			}
		}
		time.Sleep(1 * time.Second)
	}
	return st
}

// 执行upgrade 程序
func runUpgrade() bool {
	switch sysOS {
	case "windows":
		//cmd := exec.Command("cmd", "/C", "start", upgradeName)
		//_ = cmd.Start()
		cmd := exec.Command("cmd.exe", "/C start "+upgradeName, " -n "+programName)
		if err := cmd.Start(); err != nil {
			fmt.Println("运行", upgradeName, "失败:", err)
			updateLog.Println("运行", upgradeName, "失败:", err)
			return false
		}
		break
	case "linux":
		//fmt.Println("[Upgrade] deb 005-1 运行 Linux 更新....")
		err := os.Chmod(upgradeName, 0777)
		if err != nil {
			//fmt.Println("[Upgrade] Debug 修改文件权限出错！", err)
			updateLog.Println("Debug 修改文件权限出错！", err)
			return false
		}
		//cmd := exec.Command("/bin/bash", "-c", "./"+upgradeName)
		//fmt.Println("[Upgrade] deb 005-2 运行程序:", runWsUpgrade)
		cmd := exec.Command("./" + upgradeName)
		err = cmd.Start()
		if err != nil {
			//fmt.Println("[Upgrade] Debug 运行"+upgradeName+"失败! Error:", err)
			updateLog.Println("Debug 运行"+upgradeName+"失败! Error:", err)
			return false
		}
		//fmt.Println("[Upgrade] deb 005-3 运行 Linux 更新完成....")
		break
	}
	return true
}
