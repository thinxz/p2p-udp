package client

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

// 客户端定义
// --------------------
// sHost   服务器地址
// sPort   服务器端口
// cPort   本地打开端口
// srcName 本地主机标识
// dstName 远程主机标识
// srcAddr
// dstAddr
// bidiPeer
// Conn    建立连接对象
// --------------------
type Client struct {
	sHost    string `json:"sHost"`
	sPort    int    `json:"sPort"`
	cPort    int    `json:"cIP"`
	srcName  string
	dstName  string
	srcAddr  *net.UDPAddr
	dstAddr  *net.UDPAddr
	bidiPeer *net.UDPAddr
	Conn     *net.UDPConn
}

func newClient(
	sHost string, sPort int,
	cPort int, srcName string, dstName string) (c *Client, err error) {

	c = &Client{sHost: sHost, sPort: sPort, cPort: cPort, srcName: srcName, dstName: dstName}

	// 本地客户端 [固定]
	c.srcAddr = &net.UDPAddr{IP: net.IPv4zero, Port: c.cPort}
	// 服务端
	c.dstAddr = &net.UDPAddr{IP: net.ParseIP(c.sHost), Port: c.sPort}

	// 登录服务器
	c.Conn, err = net.DialUDP("udp", c.srcAddr, c.dstAddr)
	if err != nil {
		fmt.Println(err)
	}

	// 发送消息
	if _, err = c.Conn.Write([]byte("hello, I'm new peer: " + c.srcName)); err != nil {
		log.Panic(err)
	}

	fmt.Println("登录成功, 等待建立连接 ... ")

	// 读取消息
	data := make([]byte, 1024)
	n, remoteAddr, err := c.Conn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}

	// 解析对应Peer地址
	c.bidiPeer = ParseAddr(string(data[:n]))

	// 开始打洞
	fmt.Printf("执行打洞 => local:%s server:%s another: %s\n", c.srcAddr, remoteAddr, c.bidiPeer.String())
	BidiHole(c)
	return
}

func Conn(
	sHost string, sPort int,
	cPort int, srcName string, dstName string) (c *Client, err error) {

	// 建立连接
	c, err = newClient(sHost, sPort, cPort, srcName, dstName)
	if err != nil {
		fmt.Printf("conn is err : %s", err)
	}

	x := 0
	// 测试发送数据
	go func() {
		for {
			x++
			time.Sleep(10 * time.Second)
			if _, err := c.Conn.Write([]byte(" from [ " + strconv.Itoa(x) + " srcName ]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()

	// 测试接收数据
	for {
		data := make([]byte, 1024)
		n, _, err := c.Conn.ReadFromUDP(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("收到数据: %s\n", data[:n])
		}
	}
}
