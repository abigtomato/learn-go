package socket

import (
	"fmt"
	"net"
	"os"
	"testing"
	"time"
)

func TestUdpServer(t *testing.T) {
	srvAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:8003")
	if err != nil {
		fmt.Printf("net.ResolveUDPAddr Error: %v\n", err)
		return
	}
	fmt.Printf("")

	conn, err := net.ListenUDP("udp", srvAddr)
	if err != nil {
		fmt.Printf("net.ListenUDP Error: %v\n", err)
		return
	}
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)

	for {
		buf := make([]byte, 4096)
		_, cltAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("conn.ReadFromUDP Error: %v\n", err)
			return
		}
		fmt.Printf("")

		go func() {
			if _, err := conn.WriteToUDP([]byte(time.Now().String()+"\n"), cltAddr); err != nil {
				fmt.Printf("conn.WriteToUDP Error: %v\n", err)
				return
			}
		}()
	}
}

func TestUdpClient(t *testing.T) {
	conn, err := net.Dial("udp", "127.0.0.1:8003")
	if err != nil {
		fmt.Printf("net.Dial Error: %v\n", err)
		return
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Printf("conn.Read Error: %v\n", err)
				return
			}
			fmt.Printf("Server Request: %v\n", string(buf[:n]))
		}
	}()

	buf := make([]byte, 4096)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			fmt.Printf("os.Stdin.Read Error: %v\n", err)
			continue
		}
		if string(buf[:n]) == "exit\r\n" {
			break
		}
		_, _ = conn.Write(buf[:n])
	}
}
