package dynv6

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

type structDynv6Config struct {
	ApiKey     string
	Domain     string
	RecordType string
	ClientIP   string
}

var dynv6Config = &structDynv6Config{}

// 更新dns
func upDns(strUrl string) string {
	fmt.Println("正在提交更新ddns...")
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Get(strUrl)
	if err != nil {
		fmt.Println("提交更新ddns错误，Error:", err)
		return "error"
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("关闭Body错误，Error: ", err)
		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("提交更新ddns返回数据错误，Error:", err)
		return "error"
	}
	//fmt.Println(string(body))
	return string(body)
}

// dynv6.com动态域名
func upDynv6() bool {
	var apiUrl string
	//构建更新url
	if dynv6Config.RecordType == "A" {
		apiUrl = "https://ipv4.dynv6.com/api/update?"
	} else {
		apiUrl = "https://ipv6.dynv6.com/api/update?"
	}
	apiUrl = apiUrl + "hostname=" + dynv6Config.Domain + "&token=" + dynv6Config.ApiKey + "&ipv6=" + dynv6Config.ClientIP
	//尝试3次更新数据
	i := 1
	for i <= 3 {
		i++
		//更新
		res := upDns(apiUrl)
		if res == "error" {
			continue
		}
		if res == "addresses updated" || res == "addresses unchanged" {
			fmt.Println("更新ddns解析记录成功.")
			break
		} else {
			fmt.Println(res)
		}
		//等待3秒
		time.Sleep(3 * time.Second)
	}
	if i < 3 {
		return true
	} else {
		return false
	}
}

// StartDynv6 入口
func StartDynv6() {
	//dons的配置信息
	switch dynv6Config.RecordType {
	case "A": //更新ipv4地址
		fmt.Println("SET DDNS: [", time.Now().Format("2006-01-02 15:04:05"), "] ipv4->", dynv6Config.ClientIP)
		upDynv6()
		break
	case "AAAA": //获取ipv6 地址
		fmt.Println("SET DDNS: [", time.Now().Format("2006-01-02 15:04:05"), "] ipv6->", dynv6Config.ClientIP)
		//更新ipv6地址
		upDynv6()
		break
	}

}
