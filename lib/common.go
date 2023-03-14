package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const IniFile = "./ws.ini" //配置文件
var AppVer string
var HttpHeader = make(map[string]interface{}, 0)

// Ipv4Reg IPv4正则
const Ipv4Reg = `((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])`

// Ipv6Reg IPv6正则
const Ipv6Reg = `((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))`

/*
 * 自定义错误
 */
type errNo struct {
	Code    int
	Message string
}

func (c *errNo) Error() string { // 实现接口
	return c.Message
}

// MyErr 自定义错误
func MyErr(code int, Message string) error {
	return &errNo{
		Code:    code,
		Message: Message,
	}
}

// MakeIP 格式化IP地址
func MakeIP(strIP string, ipType int) string {
	if ipType == 4 {
		ipv4 := regexp.MustCompile(Ipv4Reg)
		return ipv4.FindString(strIP)
	} else {
		ipv6 := regexp.MustCompile(Ipv6Reg)
		return ipv6.FindString(strIP)
	}
}

// FileExists 文件/文件夹是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// HttpClient http请求
func HttpClient(strUrl string, strType string, strParam string) string {
	//构造http请求
	client := &http.Client{}
	req, err := http.NewRequest(strType, strUrl, strings.NewReader(strParam))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println(" Error 构造", strUrl, "请求", strType, "错误:", err)
		Logger.Println(" Error 构造", strUrl, "请求", strType, "错误:", err)
		return ""
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(" Error 发送", strUrl, "请求", strType, "错误:", err)
		Logger.Println(" Error 发送", strUrl, "请求", strType, "错误:", err)
		return ""
	}
	//返回数据
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(" Error ", strUrl, "的返回数据", body, "错误:", err)
		Logger.Println(" Error ", strUrl, "的返回数据", body, "错误:", err)
		return ""
	}
	strBody := string(body)
	//fmt.Println(strBody)
	//返回状态
	stdout := os.Stdout
	_, err = io.Copy(stdout, resp.Body)
	status := resp.StatusCode
	if status == 200 {
		return strBody
	} else {
		fmt.Println(" Error ", strUrl, "请求成功,但代码错误:", status)
		Logger.Println(" Error ", strUrl, "请求成功,但代码错误:", status)
		return "error"
	}
}

// HttpPost Post 请求
func HttpPost(url string, param string, paramType string) (_result []byte, _err error) {
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	req, err := http.NewRequest("POST", url, strings.NewReader(param))
	if err != nil {
		_err = err
	}
	//Post请求类型
	if paramType == "json" {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	//header 需要增加的内容
	if len(HttpHeader) > 0 {
		for k, v := range HttpHeader {
			req.Header.Add(k, v.(string))
		}
	}
	res, err := client.Do(req)
	if err != nil {
		_err = err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_err = err
		}
	}(res.Body)
	//返回状态
	status := res.StatusCode
	if status == 200 {
		_result, err = io.ReadAll(res.Body)
		if err != nil {
			_err = err
		}
	} else {
		errMsg, _ := io.ReadAll(res.Body)
		fmt.Println("请求:", url, "参数:", param)
		fmt.Println("请求成功，但是代码错误:", string(errMsg))
		_err = MyErr(status, "请求成功，但是代码错误:"+strconv.Itoa(status))
		Logger.Println(" Error 请求:", url, "参数:", param, "成功，但是代码错误:", status)
	}
	return _result, _err
}

// AnyJson Json通用处理
func AnyJson(byteJson []byte) map[string]interface{} {
	resMap := make(map[string]interface{}, 0)
	if err := json.Unmarshal(byteJson, &resMap); err != nil {
		fmt.Println("解析Json失败: ", err)
		Logger.Println(" Error 解析Json失败:", err)
	}
	return resMap
}

// CheckWssReturn 检查 WorkStation Server 返回信息
func CheckWssReturn(body []byte) (string, bool) {
	arrRes := make(map[string]interface{}, 0)
	err := json.Unmarshal(body, &arrRes)
	if err != nil {
		fmt.Println("解析服务器返回信息错误！")
		Logger.Println(" Error 解析服务器返回信息错误:", err)
	}
	if arrRes["code"] == "0" {
		strRes, _ := json.Marshal(arrRes["data"])
		//fmt.Println(string(strRes))
		return string(strRes), true
	} else {
		return arrRes["msg"].(string), false
	}
}

/* 字符串相关 */

// Strip 替换字符串中的字符串
func Strip(s_ string, chars_ string) string {
	s, chars := []rune(s_), []rune(chars_)
	length := len(s)
	max := len(s) - 1
	l, r := true, true //标记当左端或者右端找到正常字符后就停止继续寻找
	start, end := 0, max
	tmpEnd := 0
	charset := make(map[rune]bool) //创建字符集，也就是唯一的字符，方便后面判断是否存在
	for i := 0; i < len(chars); i++ {
		charset[chars[i]] = true
	}
	for i := 0; i < length; i++ {
		if _, exist := charset[s[i]]; l && !exist {
			start = i
			l = false
		}
		tmpEnd = max - i
		if _, exist := charset[s[tmpEnd]]; r && !exist {
			end = tmpEnd
			r = false
		}
		if !l && !r {
			break
		}
	}
	if l && r { // 如果左端和右端都没找到正常字符，那么表示该字符串没有正常字符
		return ""
	}
	return string(s[start : end+1])
}

// strInArr 字符串中含有数组中的元素 (包含不是等于)
func strInArr(str string, arrays []string) bool {
	for _, arr := range arrays {
		if strings.Contains(str, arr) {
			return true
		}
	}
	return false
}
