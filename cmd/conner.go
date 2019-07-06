package main

import (
	"flag"
	"fmt"
	"github.com/thinxz-yuan/p2p-udp/client"
	"github.com/thinxz-yuan/p2p-udp/server"
	"log"
)

var (
	d   bool
	pro string

	l string
	u string

	sHost string
	sPort int
	cPort int
)

func init() {
	flag.BoolVar(&d, "d", false, "is server")
	flag.StringVar(&pro, "pro", "UDP", "login user name")

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

		c := client.NewClient(sHost, sPort, cPort, l, u, pro)
		err := c.Connect()
		if err != nil {
			fmt.Print("客户端连接建立失败")
			return
		}

		// 执行打洞
		err = c.BidiHole()
		if err != nil {
			fmt.Print("执行打洞失败")
			return
		}

		// 检测
		c.Check()

		defer func() {
			err := c.Conn.Close()
			if err != nil {
				fmt.Printf("close conn : %s", err)
			}
		}()

		return
	}

	log.Printf("服务端启动 ...")
	// 初始化
	s := server.NewServer(pro, sPort)
	// 开启服务监听
	s.Start()
	// 等待完成, 或者服务异常退出
	<-s.CheckCountChan
}
