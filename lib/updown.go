package lib

import (
	"fmt"
	"io"
	"net"
	"os"
)

// 发送文件
func sendFile(conn net.Conn, filePath string) {
	// 只读方打开文件
	f, err := os.Open(filePath)
	defer f.Close()
	if err != nil {
		fmt.Println(" os.Open err:", err)
		return
	}
	// 循环读取文件内容,读多少发多少,原封不动
	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("文件发送完毕")
				return
			} else {
				fmt.Println(" f.Read err:", err)
				return
			}
		}
		// 将读取到的内容写到socket中,即:将本都文件写到网络中
		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Println(" conn.Write err:", err)
			return
		}
	}

}

// 接收文件
func receiveFile(conn net.Conn, fileName string) {
	// 按照文件名保存文件
	filePath := "./" + fileName
	f, err := os.Create(filePath) // 创建文件
	defer f.Close()
	if err != nil {
		fmt.Println(" os.Create err:", err)
		return
	}
	// 循环从网络中读取文件内容写入文件
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		//if n == 0{
		//	fmt.Println("接收完毕")
		//	return
		//}
		if err != nil {
			if err == io.EOF {
				fmt.Println("接收完毕")
				return
			} else {
				fmt.Println("conn.Read err:", err)
				return
			}
		}
		content := buf[:n]
		// 内容写文件,读多少写多少
		f.Write(content)
	}

}

// PutServer 发送文件服务器端
func PutServer(host string) {
	// 创建用于监听的socket
	Listener, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println(" net.Listen err:", err)
		return
	}
	defer Listener.Close()
	fmt.Println("Recv starting...")
	// 创建用于通信的socket
	conn, err := Listener.Accept()
	if err != nil {
		fmt.Println(" Listener.Accept err:", err)
		return
	}
	defer conn.Close()

	// 接收文件名
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(" conn.Read err:", err)
		return
	}
	// 获取文件名
	fileName := string(buf[:n])

	// 回复发送端 "ok"
	_, err = conn.Write([]byte("ok"))
	if err != nil {
		fmt.Println(" conn.Write err:", err)
		return
	}
	// 接收文件内容,写入到文件
	receiveFile(conn, fileName)
}

// PutClient 发送文件客户端
func PutClient(host string, filePath string) {
	// 获取文件属性
	fileIndo, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("os.Stat err:", err)
		return
	}
	fileName := fileIndo.Name()
	fileSize := fileIndo.Size()
	fmt.Println(filePath, fileName, fileSize)

	// 建立TCP连接
	conn, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	defer conn.Close()
	// 发送文件名给接收端
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}

	// 接收对端返回的确认信息
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("conn.Read err:", err)
		return
	}
	if string(buf[:n]) == "ok" {
		// 借助网络发送文件内容
		sendFile(conn, filePath)

	} else {
		fmt.Println("接收失败")
		return
	}
}

// GetClient 抓取文件客户端
func GetClient() {
	//ws类型

	//http类型

	//ftp类型

}
