package aliyun

import (
	"fmt"
	"log"
	"net/url"
)

const aliEndpoint string = "https://alidns.aliyuncs.com/"

// StructAliDnsConfig 基本参数
type StructAliDnsConfig struct {
	ApiKey     string
	ApiSecret  string
	Domain     string
	SubDomain  string
	TTL        string
	RecordType string // ipv4="A" |ipv6="AAAA"
	ClientIP   string
}

// 解析记录 返回结构
type aliDnsSubDomainRecords struct {
	TotalCount    int
	DomainRecords struct {
		Record []struct {
			DomainName string
			RecordID   string
			Value      string
		}
	}
}

// 修改/添加返回结果
type aliDnsResp struct {
	RecordID  string
	RequestID string
}

var AliDnsConfig = &StructAliDnsConfig{}
var AliDnsUpdateStatus string

func init() {

}

// 创建
func create() {
	params := url.Values{}
	params.Set("Action", "AddDomainRecord")
	params.Set("DomainName", AliDnsConfig.Domain)
	params.Set("RR", AliDnsConfig.SubDomain)
	params.Set("Type", AliDnsConfig.RecordType)
	params.Set("Value", AliDnsConfig.ClientIP)
	params.Set("TTL", AliDnsConfig.TTL)

	var result aliDnsResp
	err := aliRequest(params, &result)

	if err == nil && result.RecordID != "" {
		fmt.Printf("新增域名解析 %s 成功！IP: %s", AliDnsConfig.SubDomain, AliDnsConfig.ClientIP)
		AliDnsUpdateStatus = "success"
	} else {
		fmt.Printf("新增域名解析 %s 失败！", AliDnsConfig.SubDomain)
		AliDnsUpdateStatus = "fail"
	}
}

// 修改
func modify(record aliDnsSubDomainRecords) {
	// 相同不修改
	if len(record.DomainRecords.Record) > 0 && record.DomainRecords.Record[0].Value == AliDnsConfig.ClientIP {
		fmt.Printf("你的IP %s 没有变化, 域名 %s", AliDnsConfig.ClientIP, AliDnsConfig.SubDomain)
		return
	}

	params := url.Values{}
	params.Set("Action", "UpdateDomainRecord")
	params.Set("RR", AliDnsConfig.SubDomain)
	params.Set("RecordId", record.DomainRecords.Record[0].RecordID)
	params.Set("Type", AliDnsConfig.RecordType)
	params.Set("Value", AliDnsConfig.ClientIP)
	params.Set("TTL", AliDnsConfig.TTL)

	var result aliDnsResp
	err := aliRequest(params, &result)

	if err == nil && result.RecordID != "" {
		log.Printf("更新域名解析 %s 成功！IP: %s", AliDnsConfig.SubDomain, AliDnsConfig.ClientIP)
		AliDnsUpdateStatus = "success"
	} else {
		log.Printf("更新域名解析 %s 失败！", AliDnsConfig.SubDomain)
		AliDnsUpdateStatus = "fail"
	}
}

// RunAliDDns 更新 ali yun dns
func RunAliDDns(strIp string, recordType string) {
	AliDnsUpdateStatus = "normal"
	if recordType == "AAAA" {
		AliDnsConfig.SubDomain = AliDnsConfig.SubDomain + ".v6"
	}
	AliDnsConfig.ClientIP = strIp
	AliDnsConfig.RecordType = recordType
	var record aliDnsSubDomainRecords
	// 获取子域名信息
	params := url.Values{}
	params.Set("Action", "DescribeSubDomainRecords")
	params.Set("SubDomain", AliDnsConfig.SubDomain+"."+AliDnsConfig.Domain)
	params.Set("Type", AliDnsConfig.RecordType)
	err := aliRequest(params, &record)
	if err != nil {
		return
	}
	if record.TotalCount > 0 { // 存在
		//fmt.Println(record)
		modify(record) //更新
	} else { // 不存在
		create()
	}
}
