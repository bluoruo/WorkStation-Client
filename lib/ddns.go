package lib

import (
	aliDns "client/ddns/aliyun"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
 * DDNS 功能
 */

var DDnsServiceStatus = "none" //DDns服务状态 stopped | running

// 排除网卡
func disInterface(str string) bool {
	var arr = []string{"VMware", "vpn", "docker", "br-"}
	return strInArr(str, arr)
}

// 排除IP
func disIpAddress(str string) bool {
	var arr = []string{"192.168.80", "192.168.81", "127.0.0", "172.18.0"}
	return strInArr(str, arr)
}

// 获取本地ipv4地址
func getLocalIPv4() string {
	allNetInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("读取IPv4地址错误", err)
	}
	var ipv4 string
	for i := 0; i < len(allNetInterfaces); i++ { //所有网卡
		if disInterface(allNetInterfaces[i].Name) { //排除网卡
			continue
		}
		if (allNetInterfaces[i].Flags & net.FlagUp) != 0 { //在线网卡
			addrs, _ := allNetInterfaces[i].Addrs()
			for _, addr := range addrs {
				if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
					if ipNet.IP.To4() != nil {
						if disIpAddress(ipNet.IP.String()) { //排除IP
							continue
						}
						ipv4 = ipNet.IP.String()
						break
					}
				}
			}
		}
	}
	return MakeIP(ipv4, 4)
}

// 获取ipv4地址
func getIPv4() string {
	//ReadConfig("SYS")
	url := serverConfig.Server + "/v2/config.php"
	param := "type=ip"
	ipv4 := HttpClient(url, "POST", param)
	return MakeIP(ipv4, 4)
}

// 获取ipv6地址
func getIPv6() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("读取IPv6地址错误", err)
	}
	var ipv6 string
	for _, addr := range addrs {
		ip := regexp.MustCompile(`(\w+:){7}\w+`).FindString(addr.String())
		if strings.Count(ip, ":") == 7 {
			ipv6 = ip
			break
		}
	}
	return MakeIP(ipv6, 6)
}

// 获取域名Dns IP
func getDnsIP(domain string) string {
	var strIP string
	ips, _ := net.LookupIP(domain)
	for _, ip := range ips {
		strIP = ip.String()
		//fmt.Println(strIP)
	}
	return strIP
}

// 检查本地IP是否变化
func checkSubRecord() bool {
	ipv4 := getIPv4()
	ipv6 := getIPv6()
	if ipv4 == "" || ipv6 == "" {
		fmt.Println("IPv4地址->", ipv4)
		fmt.Println("IPv6地址->", ipv6)
		fmt.Println("无法获取ipv4或ipv6地址")
		return true
	}
	var domain = ClientBaseConfig.DdnsSubDomain
	switch ClientBaseConfig.DdnsType {
	case 1:
		domain = ClientBaseConfig.DdnsSubDomain + ".v6." + ClientDDnsConfig.MasterDomain
		dnsIP := getDnsIP(domain)
		if ipv6 == dnsIP {
			return true
		}
	case 2:
		domain = ClientBaseConfig.DdnsSubDomain + "." + ClientDDnsConfig.MasterDomain
		dnsIP := getDnsIP(domain)
		if ipv4 == dnsIP {
			return true
		}
	default:
		domain = ClientBaseConfig.DdnsSubDomain + "." + ClientDDnsConfig.MasterDomain
		dnsIP := getDnsIP(domain)
		if ipv4 == dnsIP {
			return true
		}
		domain = ClientBaseConfig.DdnsSubDomain + ".v6." + ClientDDnsConfig.MasterDomain
		dnsIP = getDnsIP(domain)
		if ipv6 == dnsIP {
			return true
		}
	}
	return false
}

// DDnsInfo 获取DDNS信息
func DDnsInfo() {
	fmt.Println("本机IP地址：", getLocalIPv4())
	ipv4 := getIPv4()
	ipv6 := getIPv6()
	iniLoad() //重载配置
	var domain = ClientBaseConfig.DdnsSubDomain
	switch ClientBaseConfig.DdnsType {
	case 1:
		domain = ClientBaseConfig.DdnsSubDomain + ".v6." + ClientDDnsConfig.MasterDomain
		fmt.Println(domain, "DNS地址->", getDnsIP(domain), "本机地址->", ipv6)
	case 2:
		domain = ClientBaseConfig.DdnsSubDomain + "." + ClientDDnsConfig.MasterDomain
		fmt.Println(domain, "DNS地址->", getDnsIP(domain), "本机地址->", ipv4)
	default:
		domain = ClientBaseConfig.DdnsSubDomain + "." + ClientDDnsConfig.MasterDomain
		fmt.Println(domain, " DNS地址->", getDnsIP(domain), "本机地址->", ipv4)
		domain = ClientBaseConfig.DdnsSubDomain + ".v6." + ClientDDnsConfig.MasterDomain
		fmt.Println(domain, " DNS地址->", getDnsIP(domain), "本机地址->", ipv6)
	}
}

// 运行 ddns aliyun 接口
func runAliyunDDns() {
	for i := 0; i < 3; i++ {
		if aliDns.AliDnsUpdateStatus == "success" {
			break
		}
		//初始化配置
		aliDns.AliDnsConfig = &aliDns.StructAliDnsConfig{
			ApiKey:    ClientDDnsConfig.ApiKey,
			ApiSecret: ClientDDnsConfig.ApiSecret,
			SubDomain: ClientBaseConfig.DdnsSubDomain,
			Domain:    ClientDDnsConfig.MasterDomain,
			TTL:       fmt.Sprintf("%v", ClientDDnsConfig.LostTime),
		}
		ipv4 := getIPv4()
		ipv6 := getIPv6()
		switch ClientBaseConfig.DdnsType {
		case 1:
			if ipv6 != "" {
				aliDns.RunAliDDns(ipv6, "AAAA")
			} else {
				fmt.Println("未获取到ipv6地址！")
			}
			break
		case 2:
			if ipv4 != "" {
				aliDns.RunAliDDns(ipv4, "A")
			} else {
				fmt.Println("未获取到ipv4地址！")
			}
			break
		default:
			if ipv4 != "" {
				aliDns.RunAliDDns(ipv4, "A")
			} else {
				fmt.Println("未获取到ipv4地址！")
			}
			if ipv6 != "" {
				aliDns.RunAliDDns(ipv6, "AAAA")
			} else {
				fmt.Println("未获取到ipv6地址！")
			}
		}
		time.Sleep(30 * time.Second)
	}
}

// 线程启动 ddns 服务
func runDDnsService() {
	switch ClientDDnsConfig.Name {
	case "aliyun":
		go runAliyunDDns()
		break
	default:
		//fmt.Println("错误的DDNS接口名")
		ddnsLog.Println(" Error 错误的DDNS接口名")
	}
}

// 检查DDns配置信息
func checkDDnsConfig() bool {
	if ClientBaseConfig.DdnsServerId == 0 ||
		ClientBaseConfig.DdnsSubDomain == "" ||
		ClientBaseConfig.DdnsType == 0 {
		fmt.Println("WorkStation Server未分配 ddns api 接口！")
		ddnsLog.Println(" Error WorkStation Server未分配 ddns api 接口！")
		fmt.Println("稍候使用 set ddns 命令再试.")
		return false
	}
	return true
}

// RunDDnsService 入口
func RunDDnsService() {
	if !checkDDnsConfig() {
		return
	}
	if DDnsServiceStatus == "running" {
		fmt.Println("DDnsService has running.")
		return
	}
	ddnsLog.Println("Start DDNS Service....")
	for {
		if DDnsServiceStatus == "stop" {
			DDnsServiceStatus = "stopped"
			break
		}
		//加载DDns配置信息
		iniLoad()
		//检查域名是否有变化
		if checkSubRecord() {
			//fmt.Println("IP地址无变化，10分钟后重试")
			ddnsLog.Println(" IP地址无变化，10分钟后重试")
		} else {
			//开始执行ddns
			DDnsServiceStatus = "running"
			runDDnsService()
		}
		//10分钟更新一次ddns解析记录
		time.Sleep(10 * time.Minute)
	}
}

// 获取 WorkStation Server DDns Api 信息
func getWwsDDnsServer() bool {
	body, err := HttpPost(serverConfig.Server+serverConfig.DDnsApi,
		"ddns="+strconv.Itoa(ClientBaseConfig.DdnsServerId),
		"normal")
	if err != nil {
		//fmt.Println("请求服务器端出错：", err)
		ddnsLog.Println(" Error 请求服务器端出错：", err)
	}
	data, st := CheckWssReturn(body)
	if st {
		err = json.Unmarshal([]byte(data), &ClientDDnsConfig)
		if err != nil {
			//fmt.Println("数据格式错误:", err)
			ddnsLog.Println(" Error 数据格式错误:", err)
		} else {
			iniWrite("DDNS") //写入配置文件
			return true
		}
	} else {
		fmt.Println("API请求错误 Error:", data)
		ddnsLog.Println(" Error API请求错误：", data)
	}
	return false
}

// UpdateWwsDDnsConfig 获取WorkStation Server DDns接口信息
func UpdateWwsDDnsConfig() bool {
	if !checkDDnsConfig() {
		return false
	}
	//获取 ddns接口信息并写入配置文件
	return getWwsDDnsServer()
}
