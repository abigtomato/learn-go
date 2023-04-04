package goroutine

import (
	"io"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

// 客户端，表示连接
type Client struct {
	// 名称
	Name string
	// 地址
	Addr string
	// 专属的信息通道
	InfoChannel chan string
}

var (
	// 全局在线连接列表
	onlineMap map[string]Client
	// 全局的消息通道
	message = make(chan string)
)

// 固定格式构造消息
func makeMsg(client Client, info string) string {
	return "[" + client.Addr + "] " + client.Name + " : " + info
}

// 写信息到指定连接
func writeToClient(client Client, conn net.Conn) {
	for msg := range client.InfoChannel {
		_, _ = conn.Write([]byte(msg + "\n"))
	}
}

// 连接保持
func connectKeep(conn net.Conn, quitChan chan bool, client Client, addr string, hasData chan bool) {
	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)

		// 正常退出
		if n == 0 {
			quitChan <- true
			log.Println("Client Exit ...")
			return
		}

		// 异常退出
		if err != nil && err != io.EOF {
			quitChan <- false
			log.Printf("conn.Read Error: %v\n", err)
			return
		}

		msg := string(buf[:n-1])
		if msg == "who" && len(msg) == 3 {
			// online
			client.InfoChannel <- "online user list: "
			for _, client := range onlineMap {
				client.InfoChannel <- client.Addr + ":" + client.Name + "\n"
			}
		} else if len(strings.Split(msg, "|")) == 2 {
			// rename
			name := strings.Split(msg, "|")[1]
			client.Name = name
			onlineMap[addr] = client
			client.InfoChannel <- "update name success"
		} else {
			// global
			message <- makeMsg(client, msg)
		}

		// 若对话顺利进行则代表用户活跃，入队标记重置select的定时器
		hasData <- true
	}
}

// 处理连接
func handlerConnect(conn net.Conn) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	// 退出标记通道
	quitChan := make(chan bool)
	// 用户活跃标记通道
	hasData := make(chan bool)

	addr := conn.RemoteAddr().String()
	client := Client{
		Name:        addr,
		Addr:        addr,
		InfoChannel: make(chan string),
	}
	onlineMap[addr] = client

	// 写信息给客户端属于耗时操作 交给新go程处理
	go writeToClient(client, conn)

	// 登录消息
	message <- makeMsg(client, "login")

	// 该go程用于保持和客户端的通信
	go connectKeep(conn, quitChan, client, addr, hasData)

	for {
		select {
		case <-quitChan:
			// 客户端退出
			delete(onlineMap, client.Addr)
			message <- makeMsg(client, "logout")
			return
		case <-hasData:
			// 若用户活跃，重置此次select的定时器
		case <-time.After(time.Second * 10):
			// 该连接超时
			delete(onlineMap, client.Addr)
			message <- makeMsg(client, "timeout")
			return
		}
	}
}

// 消息调度器
func schedule() {
	onlineMap = make(map[string]Client)

	for {
		msg := <-message
		for _, client := range onlineMap {
			client.InfoChannel <- msg
		}
	}
}

func TestChatRoom(t *testing.T) {
	// 端口监听
	listener, _ := net.Listen("tcp", "127.0.0.1:8000")
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	// 负责分发消息的go程
	go schedule()

	for {
		// 接收连接
		conn, _ := listener.Accept()

		// 每收到一个连接就开启一个专属的go程处理
		go handlerConnect(conn)
	}
}
