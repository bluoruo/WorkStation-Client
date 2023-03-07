package update

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"runtime"
)

/*
 * 文件操作
 */

// 文件/文件夹是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 文件md5值
func fileMd5Sum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		str1 := "Open Error"
		return str1, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	body, err := io.ReadAll(f)
	if err != nil {
		str2 := "io.ReadAll"
		return str2, err
	}
	strMd5 := fmt.Sprintf("%x", md5.Sum(body))
	runtime.GC()
	return strMd5, nil
}
