package server

import (
	"fmt"
	"log"
	"net"
	"time"
)

// 服务端定义
// --------------------
// Port    服务器端口
// --------------------
type Server struct {
	Port int `json:"port"`
}

// UDP 服务, 监听端口
func newServer(port int, checkCountChan chan uint64) {
	s := Server{port}

	// 启动UDP服务, 监听端口
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: s.Port})
	if err != nil {
		fmt.Println(err)
		checkCountChan <- 1
		return
	}
	log.Printf("本地地址: <%s> \n", listener.LocalAddr().String())

	// UDP 客户端连接对象
	peers := make([]net.UDPAddr, 0, 2)
	data := make([]byte, 1024)
	for {
		// 读取连接
		n, remoteAddr, err := listener.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])

		// 添加对象
		peers = append(peers, *remoteAddr)
		if len(peers) == 2 {
			//
			log.Printf("进行UDP打洞,建立 %s <--> %s 的连接\n", peers[0].String(), peers[1].String())

			// 传输连接地址
			listener.WriteToUDP([]byte(peers[1].String()), &peers[0])
			listener.WriteToUDP([]byte(peers[0].String()), &peers[1])
			time.Sleep(time.Second * 8)

			//
			log.Println("中转服务器退出,仍不影响peers间通信")

			checkCountChan <- 0
			return
		}
	}
}

// 监控运行状态
func Monitor(port int) <-chan uint64 {
	// 检查计数通道
	checkCountChan := make(chan uint64, 2)
	newServer(port, checkCountChan)
	return checkCountChan
}
