package lib

import (
	"encoding/base64"
)

/***********************
 * 自定义二次加密
***********************/

var secKey = "harry"
var secPoint = [5]int{3, 6, 9, 12, 15}

/*
 * base64 解密
 */
func base64Decode(src string) string {
	code, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		Logger.Println("解密错误: ", err)
		panic(err)
	}
	return string(code)
}

/*
 * base64 加密
 */
func base64Encode(src string) string {
	code := base64.StdEncoding.EncodeToString([]byte(src))
	return code
}

/*
 * 二次 解密
 */
func secDecode(src string) string {
	var str = src
	var point = 0
	for i := 0; i < len(secPoint); i++ {
		point = secPoint[i]
		str = str[:point-i] + str[point-i+1:]
	}
	return str
}

/*
 * 二次 加密
 */
func secEncode(src string) string {
	var str = src
	for i := 0; i < len(secPoint); i++ {
		str = str[:secPoint[i]] + secKey[i:i+1] + str[secPoint[i]:]
	}
	return str
}

// CmcheDecode 解密
func CmcheDecode(src string) string {
	str := secDecode(src)
	//fmt.Println(str)
	str = base64Decode(str)
	//fmt.Println(str)
	return str
}

// CmcheEncode 加密
func CmcheEncode(src string) string {
	str := base64Encode(src)
	//fmt.Println(str)
	str = secEncode(str)
	//fmt.Println(str)
	return str
}
