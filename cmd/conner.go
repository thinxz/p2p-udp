package main

import (
	"flag"
	"fmt"
	"github.com/thinxz-yuan/p2p-udp/client"
	"github.com/thinxz-yuan/p2p-udp/server"
)

var (
	d bool
	l string
	u string

	sHost string
	sPort int
	cPort int
)

func init() {
	flag.BoolVar(&d, "d", false, "is server")
	flag.StringVar(&l, "l", "", "login user name")
	flag.StringVar(&u, "u", "", "conn user name")

	flag.StringVar(&sHost, "sH", "47.110.253.133", "server host")
	flag.IntVar(&sPort, "sP", 9527, "server port")
	flag.IntVar(&cPort, "cP", 9250, "client port")
}

func main() {
	flag.Parse()

	if !d {
		// 客户端
		fmt.Println("启动客户端 ...")
		if l == "" {
			fmt.Println("请传入登录名 ...")
			return
		}
		c, err := client.Conn(sHost, sPort, cPort, l, u)
		if err != nil {
			fmt.Printf("conn is err : %s", err)
		}
		defer func() {
			err := c.Conn.Close()
			if err != nil {
				fmt.Printf("close conn : %s", err)
			}
		}()

		return
	}

	fmt.Println("启动服务端 ...")
	checkCountChan := server.Monitor(sPort)
	// 等待完成
	<-checkCountChan

}
