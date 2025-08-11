package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"tcp-proxy/internal/protocol"
	"tcp-proxy/pkg/log"
)

func main() {
	// 初始化日志
	log.Init("info", "")

	fmt.Println("🧪 基础3跳转发测试")
	fmt.Println("=========================")

	// 启动节点
	fmt.Println("启动转发节点...")
	go startBasicProxy(":9201", "Proxy1", "127.0.0.1:9201")
	go startBasicProxy(":9202", "Proxy2", "127.0.0.1:9202")
	go startBasicProxy(":9203", "Target", "127.0.0.1:9203")

	// 等待启动
	time.Sleep(2 * time.Second)

	// 执行测试
	fmt.Println("\n开始3跳转发测试...")
	testBasic3Hop()

	fmt.Println("\n测试完成")
}

func startBasicProxy(port, name, address string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Printf("❌ %s 启动失败: %v\n", name, err)
		return
	}

	fmt.Printf("✅ %s 启动: %s\n", name, port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleBasicConnection(conn, name, address)
	}
}

func handleBasicConnection(conn net.Conn, name, address string) {
	defer conn.Close()
	
	fmt.Printf("📦 %s 收到连接\n", name)

	// 读取数据
	buffer := make([]byte, 64*1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("❌ %s 读取失败: %v\n", name, err)
		return
	}

	data := buffer[:n]
	fmt.Printf("📨 %s 收到数据: %d字节\n", name, len(data))

	// 如果数据太短，直接返回
	if len(data) < protocol.BlockHeaderSize {
		fmt.Printf("⚠️ %s 数据太短，跳过处理\n", name)
		return
	}

	// 解析头部
	header, err := protocol.DeserializeHeader(data[:protocol.BlockHeaderSize])
	if err != nil {
		fmt.Printf("❌ %s 头部解析失败: %v\n", name, err)
		return
	}

	fmt.Printf("📋 %s 解析成功: BlockID=%d\n", name, header.BlockID)

	// 获取路径
	pathString := header.RouteInfo.TargetAddress
	fmt.Printf("🗺️ %s 路径: %s\n", name, pathString)

	// 计算下一跳
	nextHop, err := protocol.GetNextHop(pathString, address)
	if err != nil {
		fmt.Printf("❌ %s 下一跳计算失败: %v\n", name, err)
		return
	}

	if nextHop == "" {
		fmt.Printf("🎯 %s 到达终点!\n", name)
		// 提取原始数据
		originalData := data[protocol.BlockHeaderSize:]
		fmt.Printf("📄 %s 原始数据: %s\n", name, string(originalData))
		return
	}

	fmt.Printf("➡️ %s 转发到: %s\n", name, nextHop)

	// 转发
	forwardConn, err := net.DialTimeout("tcp", nextHop, 5*time.Second)
	if err != nil {
		fmt.Printf("❌ %s 连接下一跳失败: %v\n", name, err)
		return
	}
	defer forwardConn.Close()

	_, err = forwardConn.Write(data)
	if err != nil {
		fmt.Printf("❌ %s 转发失败: %v\n", name, err)
		return
	}

	fmt.Printf("✅ %s 转发成功\n", name)
}

func testBasic3Hop() {
	// 3跳路径
	path := []string{
		"127.0.0.1:9201", // Proxy1
		"127.0.0.1:9202", // Proxy2
		"127.0.0.1:9203", // Target
	}

	fmt.Printf("📍 路径: %v\n", path)

	// 测试消息
	message := "Hello 3-hop forwarding!"
	fmt.Printf("📝 消息: %s\n", message)

	// 创建数据块
	pathString := strings.Join(path, ",")
	routeInfo := protocol.RouteInfo{
		TargetAddress: pathString,
		Priority:      1,
		TTL:           64,
		Flags:         protocol.RouteFlagDirect,
		Timestamp:     time.Now().Unix(),
	}

	blockID := uint64(time.Now().UnixNano())
	dataBlock := protocol.CreateDataBlock(blockID, []byte(message), routeInfo)

	// 序列化
	serialized, err := dataBlock.Serialize()
	if err != nil {
		fmt.Printf("❌ 序列化失败: %v\n", err)
		return
	}

	fmt.Printf("📦 数据块大小: %d字节\n", len(serialized))

	// 发送到第一跳
	conn, err := net.DialTimeout("tcp", path[0], 5*time.Second)
	if err != nil {
		fmt.Printf("❌ 连接失败: %v\n", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write(serialized)
	if err != nil {
		fmt.Printf("❌ 发送失败: %v\n", err)
		return
	}

	fmt.Printf("✅ 数据已发送到第一跳\n")

	// 等待处理完成
	time.Sleep(3 * time.Second)
}
