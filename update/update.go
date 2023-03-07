package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

/*
	 -----------------------------------------
	 UPDATE 包入口文件
	 -----------------------------------------
	 Check Update Api:
	 type:		update
	 api.os:	<windows,linux,darwin,openwrt>
	 api.arch:	<amd64,x86,arm64,arm>

	 Return Json:
		{
			"ver": "0.23",
			"md5": "9375a8645b6d9bd4b49c5ae095321058",
			"os": "windows",
			"arch": "amd64",
			"down_url": "http://xx.com/download/xxx/"
		}

	 Params 调用参数:
		@param	CheckUrl	string	版本效验地址
		@param	AppVer		string	版本信息

	 CheckUrl 效验地址必须包含:
		@return	Ver			string	新版本
		@return	Md5			string	新版Md5
		@return	DownUrl		sting	新版下载地址

		CheckUrl.DownUrl:
		构造下载地址
		Windows:http://xx.com/download/xxx/<programNme>.zip (去掉.exe)
		Linux:	http://xx.com/download/xxx/<programNme>_linux_<sysArch>.zip
		OpenWrt:http://xx.com/download/xxx/<programNme>_openwrt_<sysArch>.zip
		MacOS:	http://xx.com/download/xxx/<programNme>_darwin_<sysArch>.zip


	 By Comanche Lab.
	 Date: 2023/03/02
*/

const (
	logFileUpdate = "./log/update.log"
	tmpPath       = "./new_version/" //更新缓存目录
	upgradeHost   = "127.0.0.1:7778" // upgrade 监听端口
	programHost   = "127.0.0.1:7777" // 本程序 update时 监听端口

)

// 服务器update请求信息结构
type structWscUpdateInfo struct {
	Ver     string `json:"ver"`
	Md5     string `json:"md5"`
	DownUrl string `json:"down_url"`
}

var wscUpdateInfo = &structWscUpdateInfo{}
var upgradeName = "ws_upgrade" // upgrade程序名

var (
	AppVer      string //当前运行版本
	CheckUrl    string //检查更新地址
	sysOS       string //操作系统版本
	sysArch     string //操作系统架构
	programName string //当前程序名
	DisUpdate   bool   //关闭更新
)

// 初始化
func init() {
	sysOS = runtime.GOOS
	sysArch = runtime.GOARCH
	if sysArch == "arch64" { // arch64 和 arm64通用
		sysArch = "arm64"
	}
	programName = oldProgram()
}

/*
 * 公共接口
 */

// UpgradeFinish 检查Upgrade端口 并发送 running
func UpgradeFinish() {
	//fmt.Println("[Update] Debug 日志文件:", LogFileUpdate)
	var str string
	var err error
	for i := 0; i < 5; i++ {
		//fmt.Println("[Update] debug 007 check 7778 and send running")
		//必须要没有劳斯莱斯肯定我社地
		str, err = startTcpClient("running")
		if err == nil {
			if str == "wait client" {
				break
			}
		}
		time.Sleep(1 * time.Second)
	}
}

// RunClientUpdateCom 执行更新 入口
func RunClientUpdateCom() {
	checkTempPath()     //检查更新缓存目录
	if checkVersion() { //检查版本
		if wscUpdateInfo.Ver == "" || wscUpdateInfo.Md5 == "" {
			//fmt.Println("[Update] 更新服务器无响应!")
			updateLog.Println("更新服务器无响应!")
			return
		}
		fmt.Println("[Update] 已是最新版.")
		return
	} else {
		fmt.Println("[Update] 有新版本.")
		updateLog.Println("有新版本.")
		if DisUpdate { //关闭更新
			fmt.Println("[Update] 默认关闭,版本更新！")
			updateLog.Println("默认关闭,版本更新！")
			return
		}
		//继续更新
		err := downAndCheckNewVersion() //下载并验证新版本
		if err != nil {
			//fmt.Println("[Update] 下载新版本错误:", err)
			updateLog.Println("下载新版本错误:", err)
			return
		}
		execUpgrade() //执行更新
	}
}

/*
 * 更新进程
 */
// 开始更新 入口
func execUpgrade() {
	renameByOS()            //根据操作系统命名
	updateBaseInfo()        //显示更新的基本信息
	if checkUpgradeFile() { //upgrade程序是否存在
		ServerMsg = "client update"                     //设置返回信息
		startUpdateTcpServer()                          //监听 7777
		fmt.Println("[Update] 运行 " + upgradeName + ".") //提示
		updateLog.Println("运行 " + upgradeName + ".")
		if runUpgrade() { //运行upgrade.exe
			if checkUpgradeRunStatus("start update") { //upgrade是否启动 [服务器信息 start update]
				tcpStatus = "stop"          //停止 TCP Server
				time.Sleep(3 * time.Second) //等待2s后退出
				os.Exit(0)                  //退出程序
			}
		}

	} else {
		fmt.Println("[Update] 缺少更新需要的 " + upgradeName + " 不存在！")
		updateLog.Println("缺少更新需要的 " + upgradeName + " 不存在！")
	}

}

/*
 * 包内功能
 */

// 检查新版本
func checkVersion() bool {
	//请求服务器最新版本信息
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Get(CheckUrl + "?type=update&os=" + sysOS + "&arch=" + sysArch)
	if err != nil {
		fmt.Println("[Update] Error 获取新版本的请求错误:", err)
		updateLog.Println("Error 获取新版本的请求错误:", err)
		return true
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//fmt.Println("[Update] Error 关闭连接错误错误:", err)
			updateLog.Println("Error 关闭连接错误错误:", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		//fmt.Println("[Update] Error 获取新版本的返回信息错误:", err)
		updateLog.Println("Error 获取新版本的返回信息错误:", err)
		return true
	}
	//fmt.Println(string(body))
	_ = json.Unmarshal(body, wscUpdateInfo)
	//fmt.Println(wscUpdateInfo)
	//检测是否正常返回
	if wscUpdateInfo.Ver == "" || wscUpdateInfo.Md5 == "" {
		fmt.Println("[Update] 没有能用的新版本！")
		updateLog.Println("Error 没有能用的新版本")
		return true
	}
	//对比本地版本
	if wscUpdateInfo.Ver == AppVer && wscUpdateInfo.Md5 == oldClientMd5() {
		return true
	} else {
		return false
	}
}

// 缓存目录
func checkTempPath() {
	if fileExists(tmpPath) { //不存在 新增
		//fmt.Println("[Update] 临时文件夹存在.")
		updateLog.Println("临时文件夹存在.")
	} else {
		fmt.Println("[Update] 新建，临时文件夹.")
		updateLog.Println("新建，临时文件夹.")
		err := os.Mkdir(tmpPath, 0777)
		if err != nil {
			updateLog.Println("新建临时文件夹错误:", err)
		}
	}
}

// 下载新版本
func downAndCheckNewVersion() error {
	fmt.Println("[Update] 开始下载新版本...")
	//下载新版本
	url, fileName := makeDownloadFile()
	if url == "error" {
		updateLog.Println("无有效下载地址")
		return errors.New("无有效下载地址")
	}
	err := downloadNewVersion(url, fileName)
	if err != nil {
		return err
	}
	//计算 下载的md5sum
	newMd5, err := fileMd5Sum(tmpPath + programName)
	if err != nil {
		return err
	}
	//验证 下载的md5sum和服务器的md5sum 是否一致
	if newMd5 != wscUpdateInfo.Md5 {
		updateLog.Println("下载的md5sum和服务器的不一致!")
		return errors.New("[Update] 下载的md5sum和服务器的不一致！")
	}
	fmt.Println("[Update] 新版本下载完成.")
	updateLog.Println("新版本下载完成.")
	return nil
}

// 构建下载文件地址和文件名
func makeDownloadFile() (string, string) {
	var url = wscUpdateInfo.DownUrl
	var name string
	//构建下载文件名
	if sysOS == "windows" {
		name = strings.TrimSuffix(programName, ".exe")
	} else {
		name = programName + "_" + sysOS + "_" + sysArch
	}
	//测试 不带版本号的 下载地址
	if testDownloadFile(url+name+".zip") == 200 {
		return url + name + ".zip", name
	}
	//测试 带版本号的 下载地址
	if testDownloadFile(url+wscUpdateInfo.Ver+"/"+name+".zip") == 200 {
		return url + wscUpdateInfo.Ver + "/" + name + ".zip", name
	}
	return "error", name
}

// 检测下载文件是否存在
func testDownloadFile(url string) int {
	//fmt.Println("[Update] debug 测试下载地址:", url)
	updateLog.Println("debug 测试下载地址:", url)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 404
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return 404
	}
	defer resp.Body.Close()
	status := resp.StatusCode
	fmt.Println("[Update] debug 请求状态 StatusCode:", status)
	updateLog.Println("debug 请求状态 StatusCode:", status)
	return status
}

// 下载新版本程序
func downloadNewVersion(url, name string) error {
	fmt.Println("[Update] 下载zip文件...")
	//fmt.Println("[Update] debug 从", url, "下载", tmpPath+name+".zip 文件")
	updateLog.Println("debug 从", url, "下载", tmpPath+name+".zip 文件")
	//下载
	err := downloadFile(url, tmpPath+name+".zip")
	if err != nil {
		//fmt.Println("[Update] debug Error zip下载错误:", err)
		return err
	}
	//解压
	//fmt.Println("[Update] zip下载完成,解压...")
	updateLog.Println("zip下载完成,解压...")
	err = safeUnZip(tmpPath, name)
	if err != nil { //解压失败
		//fmt.Println("[Update] debug Error zip解压错误:", err)
		updateLog.Println("debug Error zip解压错误:", err)
		return err
	}
	fmt.Println("[Update] 解压完成.")
	updateLog.Println("解压完成.")
	//不是 windows 重新命名
	if sysOS != "windows" {
		err = os.Rename(tmpPath+name, tmpPath+programName)
		if err != nil {
			//fmt.Println("[Update] Debug 重命名新程序错误！", err)
			return err
		}
		err = os.Chmod(tmpPath+programName, 0777)
		if err != nil {
			//fmt.Println("[Update] Debug 修改新程序权限错误！", err)
			return err
		}
	}
	return nil
}

// 当前程序名
func oldProgram() string {
	path, _ := os.Executable()
	_, name := filepath.Split(path)
	return name
}

// 老版本md5
func oldClientMd5() string {
	file, _ := exec.LookPath(os.Args[0])
	wscMd5, err := fileMd5Sum(file)
	if err != nil {
		//fmt.Println("[Update] Error 文件md5sum计算错误:", err)
		updateLog.Println("Error 文件md5sum计算错误:", err)
		return err.Error()
	}
	return wscMd5
}

func renameByOS() {
	if sysOS == "windows" {
		// 本程序名
		if !strings.Contains(programName, ".exe") {
			programName = programName + ".exe"
		}
		// upgrade程序名
		if !strings.Contains(upgradeName, ".exe") {
			upgradeName = upgradeName + ".exe"
		}
	} else {
		upgradeName = upgradeName + "_" + sysOS + "_" + sysArch
	}
}

// 更新基本信息
func updateBaseInfo() {
	updateLog.Println("----------------------------------------------------------------------------")
	updateLog.Println("[Update Info] 当前操作系统:", sysOS)
	updateLog.Println("[Update Info] 操作系统版本:", sysArch)
	updateLog.Println("[Update Info] 当前版本:", AppVer)
	updateLog.Println("[Update Info] 新版本:", wscUpdateInfo.Ver)
	updateLog.Println("[Update Info] 新版本下载地址:", wscUpdateInfo.DownUrl)
	updateLog.Println("[Update Info] 更新所用程序:", programName)
	updateLog.Println("----------------------------------------------------------------------------")
	updateLog.Println("[Update] 开始更新......")

	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("[Update Info] 当前操作系统:", sysOS)
	fmt.Println("[Update Info] 操作系统版本:", sysArch)
	fmt.Println("[Update Info] 当前版本:", AppVer)
	fmt.Println("[Update Info] 新版本:", wscUpdateInfo.Ver)
	fmt.Println("[Update Info] 新版本下载地址:", wscUpdateInfo.DownUrl)
	fmt.Println("[Update Info] 更新所用程序:", programName)
	fmt.Println("----------------------------------------------------------------------------")
	fmt.Println("[Update] 开始更新......")
}
