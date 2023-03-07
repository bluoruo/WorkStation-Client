package update

import (
	"fmt"
	"net"
	"time"
)

/*
 * TCP 相关
 */

var (
	ServerMsg string //Tcp Server 返回信息
	tcpStatus string //Tcp Server 运行状态
)

// 后台启动 Tcp Server
func startUpdateTcpServer() {
	if tcpStatus == "running" { //是否已经运行
		fmt.Println("[TCP:Server] has running!")
		return
	}
	go startTcpServer() //后台启动
}

// 启动 Tcp Server
func startTcpServer() {
	fmt.Println("[TCP:Server] at:", programHost)
	listen, err := net.Listen("tcp", programHost) //开启端口
	if err != nil {
		fmt.Println("[TCP:Server] listen Error:", err)
	}
	tcpStatus = "running"
	fmt.Println("[TCP:Server] started.")
	for {
		if tcpStatus == "stop" { //停止
			break
		}
		// 接受数据
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("[TCP:Server] accept Error:", err)
		}
		go acceptClientMsg(conn)
	}
	err = listen.Close() //关闭端口
	if err != nil {
		fmt.Println("[TCP:Server] Stop Error:", err)
		return
	}
	tcpStatus = "stopped"
}

// 接受客户端消息
func acceptClientMsg(conn net.Conn) {
	//fmt.Println("[TCP:Server] Debug 003 Set resMsg:", ServerMsg)
	//创建消息缓冲区
	buffer := make([]byte, 512)
	//先读取消息
	n, err := conn.Read(buffer)
	if err != nil {
		//fmt.Println("Server read msg Error:", err)
		return
	}
	clientMsg := string(buffer[0:n])
	fmt.Println("[TCP:Server] accept Msg:", clientMsg)
	//再发送消息
	conn.Write([]byte(ServerMsg))
	conn.Close()
}

/* TCP Client 相关 */

// 启动客户端
func startTcpClient(msg string) (string, error) {
	var str string
	connTimeout := 1 * time.Second
	conn, err := net.DialTimeout("tcp", upgradeHost, connTimeout)
	if err != nil {
		//fmt.Println("Port none connect,Error:", err)
		return "", err
	}
	defer conn.Close() //关闭
	// 发送消息
	fmt.Println("[TCP:Client] send to ws_upgrade Msg:", msg)
	conn.Write([]byte(msg))
	// 接受消息
	for {
		buf := [512]byte{}
		n, err := conn.Read(buf[:])
		if err != nil {
			//fmt.Println("Receive Server Error:", err)
			return "", err
		}
		str = string(buf[0:n])
		if str != "" {
			break
		}
		time.Sleep(1 * time.Second)
	}

	//fmt.Println("Client Exit")
	return str, nil
}

/* 端口扫描 */

// 端口扫描
func scanPort(host string) error {
	connTimeout := 1 * time.Second
	//conn, err := net.Dial("tcp", host+":"+port)
	conn, err := net.DialTimeout("tcp", host, connTimeout)
	if err != nil {
		fmt.Println("Port none connect,Error:", err)
		return err
	}

	//关闭
	err = conn.Close()
	if err != nil {
		fmt.Println("Stop tcp client Error:", err)
		return err
	}
	return nil
}
