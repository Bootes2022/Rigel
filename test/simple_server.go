package main

import (
	"fmt"
	"io"
	"net"
	"log"
)

func main() {
	// 启动简单的测试服务器
	listener, err := net.Listen("tcp", "127.0.0.1:8081")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer listener.Close()

	fmt.Println("Test server listening on 127.0.0.1:8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept error: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	
	fmt.Printf("New connection from %s\n", conn.RemoteAddr())
	
	// 读取数据并回显
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("Read error: %v\n", err)
			}
			break
		}
		
		data := buffer[:n]
		fmt.Printf("Received: %s", string(data))
		
		// 回显数据
		response := fmt.Sprintf("Echo: %s", string(data))
		conn.Write([]byte(response))
	}
	
	fmt.Printf("Connection from %s closed\n", conn.RemoteAddr())
}
