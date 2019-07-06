package client

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	reuse "github.com/thinxz-yuan/go-reuseport"
)

// 打洞业务
func (c *Client) BidiHole() (err error) {
	// 执行打洞
	log.Printf("执行打洞 => local:%s server:%s another: %s\n", c.srcAddr, c.dstAddr, c.bidiPeer)

	// 必须先关闭然后再建立连接
	err = c.Conn.Close()
	if err != nil {
		fmt.Printf("must close UDPConn before BidiHole : %s", err)
		return
	}

	// 建立点对点连接 [保存连接并返回]
	c.Conn, err = reuse.Dial(c.network, c.srcAddr, c.bidiPeer)
	if err != nil {
		log.Printf("<%s> BidiHole 失败 => <%s -> %s>  , %s  ... \n", c.network, c.srcAddr, c.bidiPeer, err)
		return
	}

	// 向另一个peer发送一条udp消息
	// (对方peer的nat设备会丢弃该消息,非法来源),用意是在自身的nat设备打开一条可进入的通道,这样对方peer就可以发过来udp消息
	if _, err = c.Conn.Write([]byte("打洞消息")); err != nil {
		log.Println("send handshake:", err)
	}
	return
}

// 解析节点
// --------------------
// addr 节点地址结构 [0.0.0.0:xx]
// --------------------
func ParseAddr(addr string) *net.UDPAddr {
	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])

	return &net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}
