package aliyun

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

// 统一请求接口
func aliRequest(params url.Values, result interface{}) (err error) {
	//构建签名
	aliyunSigner(AliDnsConfig.ApiKey, AliDnsConfig.ApiSecret, &params)
	//构造请求参数
	req, err := http.NewRequest("GET", aliEndpoint, bytes.NewBuffer(nil))
	req.URL.RawQuery = params.Encode()
	if err != nil {
		log.Println("http.NewRequest失败. Error: ", err)
		return
	}
	client := &http.Client{
		Timeout: time.Second * 20,
	}
	resp, err := client.Do(req)
	//处理返回
	err = getHTTPResponse(resp, aliEndpoint, err, result)

	return
}

// 处理返回 序列化的json
func getHTTPResponse(resp *http.Response, url string, err error, result interface{}) error {
	body, err := getHTTPResponseOrg(resp, url, err)

	if err == nil {
		// log.Println(string(body))
		err = json.Unmarshal(body, &result)

		if err != nil {
			log.Printf("请求接口%s解析json结果失败! ERROR: %s\n", url, err)
		}
	}

	return err

}

// 处理返回 byte
func getHTTPResponseOrg(resp *http.Response, url string, err error) ([]byte, error) {
	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Printf("请求接口%s失败! ERROR: %s\n", url, err)
	}

	// 300及以上状态码都算异常
	if resp.StatusCode >= 300 {
		errMsg := fmt.Sprintf("请求接口 %s 失败! 返回内容: %s ,返回状态码: %d\n", url, string(body), resp.StatusCode)
		log.Println(errMsg)
		err = fmt.Errorf(errMsg)
	}

	return body, err
}
