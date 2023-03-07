package update

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
)

/*
 * ZIP 解压
 */

// 安全的解压文件 入口
func safeUnZip(path, name string) error {
	//开始解压
	err := unZip(path, name+".zip")
	if err != nil {
		return err
	}
	//解压后的处理
	name = path + name
	var unName string
	if sysOS == "windows" {
		unName = name + ".exe"
	} else {
		unName = name
	}
	//解压后文件是否存在
	if fileExists(unName) { //不存在 报错
		//存在 删除zip文件
		err = os.Remove(name + ".zip")
		if err != nil {
			return err
		}
	} else {
		err = errors.New("文件不存在！Error:" + unName)
		return err
	}
	return nil
}

// 解压zip文件
func unZip(path, zipFile string) error {
	archive, err := zip.OpenReader(path + zipFile)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(path, f.Name)
		//filePath := f.Name
		//fmt.Println("unzipping file:", filePath)

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return err
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return err
		}

		dstFile.Close()
		fileInArchive.Close()
	}
	return nil
}
