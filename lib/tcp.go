package lib

import (
	"fmt"
	"net"
	"os"
	"time"
)

// 接收连接
func handleTcpConnection(conn net.Conn) {
	buffer := make([]byte, 2048) //建立一个slice
	for {
		n, err := conn.Read(buffer) //读取客户端传来的内容
		if err != nil {
			fmt.Println(err)
			return //当远程客户端连接发生错误（断开）后，终止此协程。
		}
		fmt.Println(conn.RemoteAddr().String(), "receive data string:\n", string(buffer[:n]))

		//返回给客户端的信息
		strTemp := "CofoxServer got msg \"" + string(buffer[:n]) + "\" at " + time.Now().String()
		conn.Write([]byte(strTemp))
	}
}

// 发送信息
func senderTcpServer(conn net.Conn) {
	words := "Hello Server!"
	conn.Write([]byte(words))
	fmt.Println("send over")

	//接收服务端反馈
	buffer := make([]byte, 2048)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(conn.RemoteAddr().String(), "waiting server back msg error: ", err)
		return
	}
	fmt.Println(conn.RemoteAddr().String(), "receive server back msg: ", string(buffer[:n]))

}

// ConnectTCPServer 连接服务器端
func ConnectTCPServer(host string) {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", host)
	if err != nil {
		fmt.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	fmt.Println("connection success")
	senderTcpServer(conn)
}

// ListenTCPServer  监听服务器端
func ListenTCPServer(host string) {
	//建立socket端口监听
	netListen, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Println("Listen to Server Error: ", err)
	}
	defer netListen.Close()
	fmt.Println("Waiting for clients ...")

	//等待客户端访问
	for {
		conn, err := netListen.Accept() //监听接收
		if err != nil {
			continue //如果发生错误，继续下一个循环。
		}
		fmt.Println(conn.RemoteAddr().String(), "tcp connect success") //tcp连接成功
		go handleTcpConnection(conn)
	}
}
