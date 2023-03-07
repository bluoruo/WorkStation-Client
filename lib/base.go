package lib

import (
	"encoding/json"
	"fmt"
)

type structBaseName struct {
	Name     string `json:"name"`
	HostId   string `json:"host_id"`
	HostName string `json:"host_name"`
}

var structBase = &structBaseName{}

// 检查 识别名
func checkBaseName(str string) bool {
	structBase.Name = str
	structBase.HostId = baseInfo.HostId
	req, _ := json.Marshal(structBase)
	body, err := HttpPost(serverConfig.Server+serverConfig.BaseApi+"?type=checkname",
		string(req),
		"json")
	if err != nil {
		Logger.Println(" [Base]Error 请求服务器端出错：", err)
	}
	data, st := CheckWssReturn(body)
	if st {
		res := AnyJson([]byte(data))
		if res["Status"] != "error" {
			//fmt.Println(res["Name"])
			//写入配置文件
			iniKeyValue("BASE", "Name", res["Name"].(string))                   //Name
			iniKeyValue("BASE", "HostId", fmt.Sprintf("%v", structBase.HostId)) //HostId
			iniKeyValue("BASE", "HostName", baseInfo.HostName)                  //HostName
			return true
		}
	} else {
		fmt.Println("API请求错误 Error:", data)
		Logger.Println(" [Base]Error API请求错误:", data)
	}
	return false
}

// SendToWssBase 获取WorkStation Server Base信息
func SendToWssBase(param interface{}) bool {
	req, _ := json.Marshal(param)
	body, err := HttpPost(serverConfig.Server+serverConfig.BaseApi+"?type=getinfo",
		string(req),
		"json")
	if err != nil {
		Logger.Println(" [Base]Error 请求服务器端出错", err)
	}
	data, st := CheckWssReturn(body)
	if st {
		err = json.Unmarshal([]byte(data), &ClientBaseConfig)
		if err != nil {
			Logger.Println(" [Base]Error 数据格式错误:", err)
		} else {
			iniWrite("BASE") //写入配置文件
			return true
		}
	} else {
		fmt.Println("API请求错误 Error:", data)
		Logger.Println(" [Base]Error API请求错误:", err)
	}
	return false
}

// GetWwsBaseConfig 获取 WorkStation Server Base信息
func GetWwsBaseConfig() bool {
	Logger.Println(" [Base]开始获取基本信息.")
	iniLoad()
	structBase.Name = ClientBaseConfig.Name
	structBase.HostId = ClientBaseConfig.HostId
	structBase.HostName = ClientBaseConfig.HostName
	//fmt.Println(structBase)
	return SendToWssBase(structBase)
}
