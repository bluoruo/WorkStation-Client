package lib

import (
	"encoding/json"
	"fmt"
	"github.com/go-ini/ini"
	"os"
	"strconv"
	"time"
)

/*
 * 配置文件 功能
 */

// Config 系统接口配置
type structServerConfig struct {
	Server     string
	BaseApi    string
	DDnsApi    string
	RepoApi    string
	StatusApi  string
	ComApi     string
	UpdateTime int64
}

// Base 基础配置
type structBaseConfig struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	HostId        string `json:"host_id"`
	HostName      string `json:"host_name"`
	DdnsServerId  int    `json:"ddns_server_id"`
	DdnsSubDomain string `json:"ddns_sub_domain"`
	DdnsType      int    `json:"ddns_type"`
	Services      string `json:"services"`
	Remark        string `json:"remark"`
}

// DDns 动态域名配置
type structDDnsConfig struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Account      string `json:"account"`
	Password     string `json:"password"`
	ApiKey       string `json:"api_key"`
	ApiSecret    string `json:"api_secret"`
	ApiToken     string `json:"api_token"`
	MasterDomain string `json:"master_domain"`
	LostTime     int    `json:"lost_time"`
	Url          string `json:"url"`
}

// Repo 上报信息配置
type structRepoConfig struct {
	Name     string
	Interval string
}

// Com 系统命令配置
type structComConfig struct {
	Interval string
}

var forceUpdate = false
var ConfURL = "https://ws.ofbao.com/v2/config.php" //默认接口
var serverConfig = &structServerConfig{}           //服务器端配置信息
var ClientBaseConfig = &structBaseConfig{}
var ClientDDnsConfig = &structDDnsConfig{}

// 读取配置信息
func iniLoad() bool {
	// 加载配置文件
	cfg, err := ini.Load(IniFile)
	if err != nil {
		fmt.Println(" [Config]error 读取配置文件错误:", err)
		Logger.Println(" error 读取配置文件错误:", err)
		return false
	}
	//获取Server信息
	err = cfg.Section("SERVER").MapTo(serverConfig)
	if err != nil {
		fmt.Println(" [Config]error 读取[SERVER]分区错误:", err)
		return false
	}

	//获取BASE信息
	err = cfg.Section("BASE").MapTo(ClientBaseConfig)
	if err != nil {
		fmt.Println(" [Config]error to [BASE]分区错误:", err)
		return false
	}

	//获取DDNS信息
	err = cfg.Section("DDNS").MapTo(ClientDDnsConfig)
	if err != nil {
		fmt.Println(" [Config]error to [DDNS]分区错误:", err)
		return false
	}

	return true
}

// 创建配置信息
func iniCreate() bool {
	cfg := ini.Empty()
	err := cfg.Section("SERVER").ReflectFrom(&serverConfig)
	if err != nil {
		fmt.Println(" [Config]error 创建[SERVER]分区错误:", err)
		return false
	}
	err = cfg.Section("BASE").ReflectFrom(&ClientBaseConfig)
	if err != nil {
		fmt.Println(" [Config]error 创建[BASE]分区错误:", err)
		return false
	}
	err = cfg.Section("DDNS").ReflectFrom(&ClientDDnsConfig)
	if err != nil {
		fmt.Println(" [Config]error 创建[DDNS]分区错误:", err)
		return false
	}
	// 保存到文件
	err = cfg.SaveTo(IniFile)
	if err != nil {
		fmt.Println(" [Config]error 创建配置文件错误: ", err)
		return false
	}
	return true
}

// 分区写入配置
func iniWrite(section string) bool {
	cfg, _ := ini.Load(IniFile)
	switch section {
	case "SERVER": // 添加sys配置信息
		err := cfg.Section("SERVER").ReflectFrom(&serverConfig)
		if err != nil {
			fmt.Println(" error 写入[SERVER]分区错误:", err)
			return false
		}
		break
	case "BASE": // 添加base配置信息
		err := cfg.Section("BASE").ReflectFrom(&ClientBaseConfig)
		if err != nil {
			fmt.Println(" error 写入[BASE]分区错误:", err)
			return false
		}
		break
	case "DDNS": // 添加ddns配置信息
		err := cfg.Section("DDNS").ReflectFrom(&ClientDDnsConfig)
		if err != nil {
			fmt.Println(" 写入[DDNS]分区错误:", err)
			return false
		}
		break
	default:
		return false
	}
	// 保存到文件
	err := cfg.SaveTo(IniFile)
	if err != nil {
		fmt.Println(" error 写入配置文件出错:", err)
		return false
	}
	return true
}

// 写入元素
func iniKeyValue(section string, key string, value string) {
	cfg, _ := ini.Load(IniFile)
	cfg.Section(section).Key(key).SetValue(value)
	// 保存到文件
	err := cfg.SaveTo(IniFile)
	if err != nil {
		Logger.Println(" [Config]error 元素写入配置文件出错:", err)
	}
	return
}

// 更新配置信息
func iniUpdateTime(sec string) bool {
	cfg := ini.Empty()
	_, err := cfg.Section(sec).NewKey("time", strconv.FormatInt(time.Now().Unix(), 10))
	if err != nil {
		fmt.Println(" [Config] 更新update时间[", sec, "]错误:", err)
	}
	return true
}

// getServerConfig 获取配置文件
func getServerConfig() bool {
	postParam := "type=config"
	strBody := HttpClient(ConfURL, "POST", postParam)
	//fmt.Print(strBody)
	if strBody == "error" {
		Logger.Println(" [Config]抓取不到配置信息！")
		return false
	} else {
		Logger.Println(" [Config] Debug 原始:", strBody)
		strBody = CmcheDecode(strBody)
		Logger.Println(" [Config] Debug 解密后", strBody)
		//返回参数装入 serverConfig
		err := json.Unmarshal([]byte(strBody), &serverConfig)
		if err != nil {
			fmt.Println(err)
			return false
		}
		serverConfig.UpdateTime = time.Now().Unix()
		return true
	}

}

// ServerConfig 服务器接口
func ServerConfig() bool {
	if !FileExists(IniFile) { //不存在 创建配置文件
		Logger.Println(" [Config]配置文件不存在,新建")
		_, err := os.OpenFile(IniFile, os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			Logger.Println(" [Config]error 创建配置文件错误:", err)
			//Logger.Println("创建配置文件错误:", err)
			return false
		}
		iniCreate() //写入空白信息
	}
	iniLoad() //读取配置文件
	now := time.Now().Unix()
	// 1小时内(不再更新配置信息) 并没有强制更新
	if (now-serverConfig.UpdateTime) < 3600 && forceUpdate == false {
		Logger.Println(" [Config]配置文件无需更新！")
		return true
	}
	// 获取新配置
	Logger.Println(" [Config]抓取服务器配置信息...")
	if getServerConfig() {
		if iniWrite("SERVER") { //写入配置文件
			Logger.Println(" [Config]更新服务器配置信息完成")
			return true
		} else {
			Logger.Println(" [Config]error 更新服务器配置信息错误！")
		}
	}
	return false
}

// 更新ServerConfig
func updateServerConfig() bool {
	if getServerConfig() {
		if iniWrite("SERVER") { //写入配置文件
			Logger.Println(" [Config]更新ServerConfig完成")
			return true
		} else {
			Logger.Println(" [Config]error 更新ServerConfig错误！")
		}
	}
	return false
}
