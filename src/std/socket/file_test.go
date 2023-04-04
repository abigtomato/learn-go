package socket

import (
	"fmt"
	"io"
	"net"
	"os"
	"testing"
)

func receiveFile(conn net.Conn, fileName string) (err error) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("os.Create Error: %v\n", err)
		return
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf := make([]byte, 4096)
	for {
		n, err := conn.Read(buf)
		if n == 0 {
			if err != nil && err == io.EOF {
				fmt.Println("客户端断开连接")
			} else {
				fmt.Printf("conn.Read Error: %v\n", err)
			}
			return
		}

		_, err = f.Write(buf[:n])
		if err != nil {
			fmt.Printf("f.Write Error: %v\n", err)
			return
		}
	}
}

func TestFileReceiver(t *testing.T) {
	listen, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("net.Listen Error: %v\n", err)
		return
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	conn, err := listen.Accept()
	if err != nil {
		fmt.Printf("listen.Accept Error: %v\n", err)
		return
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("conn.Read Error: %v\n", err)
		return
	}

	_, err = conn.Write([]byte("ok"))
	if err != nil {
		fmt.Printf("conn.Write Error: %v\n", err)
		return
	}

	if err := receiveFile(conn, string(buf[:n])); err != nil {
		fmt.Printf("recvFile Error: %v\n", err)
		return
	}
}

func sendFile(conn net.Conn, file string) error {
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("os.Open Error: %v\n", err)
		return err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	buf := make([]byte, 4096)
	for {
		n, err := f.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("%v\n", err)
			} else {
				fmt.Printf("f.Read Error: %v\n", err)
			}
			return err
		}

		_, err = conn.Write(buf[:n])
		if err != nil {
			fmt.Print(err)
			return err
		}
	}
}

func TestFileSender(t *testing.T) {
	list := os.Args
	if len(list) != 2 {
		fmt.Println("格式错误")
		return
	}

	filePath := list[1]
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("os.Stat Error: %v\n", err)
		return
	}
	fileName := fileInfo.Name()

	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("net.Dial Error: %v\n", err)
		return
	}
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Printf("conn.Write Error: %v\n", err)
		return
	}

	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("conn.Read Error: %v\n", err)
		return
	}

	if "ok" == string(buf[:n]) {
		err := sendFile(conn, filePath)
		if err != nil {
			return
		}
	} else {
		fmt.Printf("Ok Error: %v\n", err)
		return
	}
}
