package server

import (
	"fmt"
	"log"
	"net"
	"time"

	reuse "github.com/thinxz-yuan/go-reuseport"
)

type UDPAddr struct {
	// 登录名
	name string
	// UDP
	uConn *net.UDPAddr
}

type TCPAddr struct {
	name string
	// TCP
	conn net.Conn
}

// 服务端定义
// --------------------
// Port    服务器端口
// --------------------
type Server struct {
	Host        string `json:"host"`
	Port        int    `json:"port"`
	network     string // 当前使用的协议 [TCP UDP]
	addrUDP     net.Addr
	udpListener *net.UDPConn

	addrTCP     net.Addr
	tcpListener net.Listener
	allAddr     []UDPAddr // 所有登录的数据
	allConn     []TCPAddr // 所有TCP登录连接
}

// 初始化服务参数
func Init(network string, port int) *Server {
	return &Server{
		Port: port, Host: net.IPv4zero.String(),
		network: network,
		addrUDP: &net.UDPAddr{IP: net.IPv4zero, Port: port},
		addrTCP: &net.TCPAddr{IP: net.IPv4zero, Port: port},
	}
}

// UDP 服务, 监听端口
func (s *Server) newServer(checkCountChan chan uint64) (err error) {
	// 启动UDP服务, 监听端口
	if s.network == "UDP" {
		s.udpListener, err = net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: s.Port})
		if err != nil {
			fmt.Println(err)
			// 错误结束
			checkCountChan <- 1
			return
		}

		r := make(chan int)
		s.allAddr = make([]UDPAddr, 0, 10)
		// 获取数据
		go func() {
			for {
				data := make([]byte, 1024)
				// 读取连接
				n, remoteAddr, err := s.udpListener.ReadFromUDP(data)
				if err != nil {
					fmt.Printf("error during read: %s", err)
				}
				log.Printf("<%s> %s\n", remoteAddr.String(), data[:n])

				// 添加对象
				s.allAddr = append(s.allAddr, UDPAddr{name: "", uConn: remoteAddr})
				// 通知写入数据
				r <- 1
			}
		}()

		go func() {
			for {
				<-r
				log.Printf("is client conn -> %d .. \n", len(s.allAddr))
				if len(s.allAddr) == 2 {
					//
					log.Printf("进行UDP打洞,建立 %s <--> %s 的连接\n", s.allAddr[0].uConn.String(), s.allAddr[1].uConn.String())
					// 传输连接地址
					_, _ = s.udpListener.WriteToUDP([]byte(s.allAddr[0].uConn.String()), s.allAddr[1].uConn)
					_, _ = s.udpListener.WriteToUDP([]byte(s.allAddr[1].uConn.String()), s.allAddr[0].uConn)
					time.Sleep(time.Second * 8)
					log.Println("中转服务器退出,仍不影响peers间通信")
					log.Printf("本地地址: <%s> \n", s.addrUDP.String())
					// 通知主线程结束
					checkCountChan <- 0
				}
			}
		}()
	} else if "TCP" == s.network {
		// 启动TCP服务, 监听端口
		s.tcpListener, err = reuse.Listen("tcp", s.addrTCP.String())
		if err != nil {
			fmt.Println(err)
			// 错误结束
			checkCountChan <- 1
			return
		}

		log.Printf("TCP 本地地址: <%s> \n", s.addrTCP.String())

		r := make(chan int)
		s.allConn = make([]TCPAddr, 0, 10)
		// 获取数据
		go func() {
			for {
				// 读取连接
				// 读取连接
				conn, err := s.tcpListener.Accept()
				if err != nil {
					fmt.Printf("error during read: %s", err)
				}
				log.Printf("<%s> %s\n", conn.RemoteAddr(), conn.LocalAddr())

				// 添加对象
				s.allConn = append(s.allConn, TCPAddr{name: "", conn: conn})

				// 通知写入数据
				r <- 1
			}
		}()

		go func() {
			for {
				<-r
				log.Printf("is client conn => %d .. \n", len(s.allConn))
				if len(s.allConn) == 2 {
					log.Printf("进行TCP打洞,建立 %s <--> %s 的连接\n", s.allConn[0].conn.RemoteAddr().String(), s.allConn[1].conn.RemoteAddr().String())

					// 传输连接地址
					_, _ = s.allConn[0].conn.Write([]byte(s.allConn[1].conn.RemoteAddr().String()))
					_, _ = s.allConn[1].conn.Write([]byte(s.allConn[0].conn.RemoteAddr().String()))
					time.Sleep(time.Second * 8)

					//
					log.Println("中转服务器退出,仍不影响peers间通信")

					checkCountChan <- 0
				}
			}
		}()
	}
	return
}

// 监控运行状态
func (s *Server) Monitor() <-chan uint64 {
	// 检查计数通道
	checkCountChan := make(chan uint64, 2)
	_ = s.newServer(checkCountChan)
	return checkCountChan
}
