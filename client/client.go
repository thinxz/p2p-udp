package client

import (
	"errors"
	"log"
	"net"
	"strconv"
	"time"

	reuse "github.com/thinxz-yuan/go-reuseport"
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
	cPort   int `json:"cIP"`
	network string

	srcName string
	dstName string

	srcAddr string
	dstAddr string

	bidiPeer string
	Conn     net.Conn
}

func NewClient(sHost string, sPort, cPort int, srcName, dstName, network string) (c *Client) {
	return &Client{
		cPort:   cPort,
		network: network,
		srcName: srcName,
		dstName: dstName,
		srcAddr: net.IPv4zero.String() + ":" + strconv.Itoa(cPort),
		dstAddr: sHost + ":" + strconv.Itoa(sPort),
	}
}

func (c *Client) Connect() (err error) {
	if "UDP" == c.network || "udp" == c.network {
		c.network = "udp"
	} else if "TCP" == c.network || "tcp" == c.network {
		c.network = "tcp"
	} else {
		log.Printf("请输入正确协议标识 <%s> -> <UDP/udp 、TCP/tcp> ... \n", c.network)
		return errors.New("not found network protocol")
	}

	c.Conn, err = reuse.Dial(c.network, c.srcAddr, c.dstAddr)

	if err != nil {
		log.Printf("<%s> 端口服务 <%s> 打开失败 ... \n", c.network, c.dstAddr)
		return
	}

	// 发送消息
	if _, err = c.Conn.Write([]byte("hello, I'm new peer: " + c.srcName)); err != nil {
		return
	}
	log.Printf("<%s> 登录 <%s> 成功, 等待建立连接 ... \n", c.network, c.dstAddr)

	// 读取消息
	data := make([]byte, 1024)
	n, err := c.Conn.Read(data)
	if err != nil {
		log.Printf("error during read: %s", err)
	}

	// 解析对应Peer地址
	c.bidiPeer = ParseAddr(string(data[:n])).String()

	log.Printf("连接建立成功, 获取目标对应地址 => <%s>... \n", c.bidiPeer)

	return
}

func (c *Client) Check() {
	x := 0
	// 测试发送数据
	go func() {
		for {
			x++
			time.Sleep(3 * time.Second)
			if _, err := c.Conn.Write([]byte(" from [ " + strconv.Itoa(x) + " srcName ]")); err != nil {
				log.Println("send msg fail", err)
			}
		}
	}()

	// 测试接收数据
	for {
		data := make([]byte, 1024)
		n, err := c.Conn.Read(data)
		if err != nil {
			log.Printf("error during read: %s\n", err)
		} else {
			log.Printf("收到数据: %s\n", data[:n])
		}
	}
}
