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

	fmt.Println("🧪 简化版本地3跳转发验证")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("使用简单TCP连接验证基础转发功能")
	fmt.Println()

	// 启动简单的转发节点
	go startSimpleForwarder("127.0.0.1:9101", "Proxy1")
	go startSimpleForwarder("127.0.0.1:9102", "Proxy2")
	go startSimpleTarget("127.0.0.1:9103", "Target")

	// 等待节点启动
	time.Sleep(1 * time.Second)

	// 执行简化验证测试
	runSimpleVerificationTests()

	fmt.Println("\n🎉 简化版3跳转发验证完成！")
}

func startSimpleForwarder(address, name string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("❌ %s 启动失败: %v\n", name, err)
		return
	}

	fmt.Printf("🔗 %s 启动成功: %s\n", name, address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleSimpleForwarding(conn, name, address)
	}
}

func startSimpleTarget(address, name string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("❌ %s 启动失败: %v\n", name, err)
		return
	}

	fmt.Printf("🎯 %s 启动成功: %s\n", name, address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleTargetReceive(conn, name)
	}
}

func handleSimpleForwarding(clientConn net.Conn, proxyName, proxyAddress string) {
	defer clientConn.Close()

	fmt.Printf("📦 %s 接收连接: %s\n", proxyName, clientConn.RemoteAddr())

	// 读取完整数据
	buffer := make([]byte, 1024*1024) // 1MB缓冲区
	totalRead := 0
	
	for {
		n, err := clientConn.Read(buffer[totalRead:])
		if err != nil {
			break
		}
		totalRead += n
		
		// 检查是否读取了完整的数据块头部
		if totalRead >= protocol.BlockHeaderSize {
			// 解析头部获取数据块大小
			header, err := protocol.DeserializeHeader(buffer[:protocol.BlockHeaderSize])
			if err == nil {
				expectedSize := protocol.BlockHeaderSize + int(header.BlockSize)
				if totalRead >= expectedSize {
					// 读取完整，处理数据块
					data := buffer[:expectedSize]
					processSimpleDataBlock(data, proxyName, proxyAddress)
					break
				}
			}
		}
		
		// 防止无限读取
		if totalRead > 1024*1024 {
			break
		}
	}
}

func processSimpleDataBlock(data []byte, proxyName, proxyAddress string) {
	// 解析头部
	header, err := protocol.DeserializeHeader(data[:protocol.BlockHeaderSize])
	if err != nil {
		fmt.Printf("❌ %s 头部解析失败: %v\n", proxyName, err)
		return
	}

	// 提取数据部分
	blockData := data[protocol.BlockHeaderSize:]
	fmt.Printf("📨 %s 解析数据块: BlockID=%d, 数据大小=%d字节\n", 
		proxyName, header.BlockID, len(blockData))

	// 解析路径信息
	pathString := header.RouteInfo.TargetAddress
	path := protocol.ParsePath(pathString)
	
	// 获取下一跳
	nextHop, err := protocol.GetNextHop(pathString, proxyAddress)
	if err != nil {
		fmt.Printf("❌ %s 获取下一跳失败: %v\n", proxyName, err)
		return
	}

	if nextHop == "" {
		fmt.Printf("🎯 %s 已到达最终目标\n", proxyName)
		return
	}

	fmt.Printf("➡️  %s 转发到下一跳: %s\n", proxyName, nextHop)

	// 转发到下一跳
	if err := forwardToNextHopSimple(data, nextHop, proxyName); err != nil {
		fmt.Printf("❌ %s 转发失败: %v\n", proxyName, err)
		return
	}

	fmt.Printf("✅ %s 转发成功\n", proxyName)
}

func forwardToNextHopSimple(data []byte, nextHop, proxyName string) error {
	// 连接到下一跳
	conn, err := net.DialTimeout("tcp", nextHop, 5*time.Second)
	if err != nil {
		return fmt.Errorf("连接下一跳失败: %v", err)
	}
	defer conn.Close()

	// 转发完整数据块
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("转发数据失败: %v", err)
	}

	return nil
}

func handleTargetReceive(conn net.Conn, targetName string) {
	defer conn.Close()

	fmt.Printf("🎯 %s 接收连接: %s\n", targetName, conn.RemoteAddr())

	// 读取数据
	buffer := make([]byte, 1024*1024)
	totalRead := 0
	
	for {
		n, err := conn.Read(buffer[totalRead:])
		if err != nil {
			break
		}
		totalRead += n
		
		if totalRead >= protocol.BlockHeaderSize {
			header, err := protocol.DeserializeHeader(buffer[:protocol.BlockHeaderSize])
			if err == nil {
				expectedSize := protocol.BlockHeaderSize + int(header.BlockSize)
				if totalRead >= expectedSize {
					// 处理完整数据块
					originalData := buffer[protocol.BlockHeaderSize:expectedSize]
					fmt.Printf("🎉 %s 接收到完整数据: BlockID=%d, 原始数据=%s\n", 
						targetName, header.BlockID, string(originalData[:min(50, len(originalData))]))
					break
				}
			}
		}
	}
}

func runSimpleVerificationTests() {
	fmt.Println("\n🧪 开始简化版3跳转发测试")
	fmt.Println(strings.Repeat("-", 40))

	// 定义3跳转发路径
	threehopPath := []string{
		"127.0.0.1:9101", // Proxy1
		"127.0.0.1:9102", // Proxy2
		"127.0.0.1:9103", // Target
	}

	fmt.Printf("📍 3跳转发路径: %v\n", threehopPath)

	// 测试消息
	testMessages := []string{
		"Test message 1: Simple 3-hop forwarding",
		"Test message 2: Verifying path-based routing",
		"Test message 3: End-to-end data integrity check",
	}

	for i, message := range testMessages {
		fmt.Printf("\n📤 发送测试消息 %d: %s\n", i+1, message)
		
		if err := sendSimpleTestMessage(message, threehopPath); err != nil {
			fmt.Printf("❌ 测试消息 %d 发送失败: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ 测试消息 %d 发送成功\n", i+1)
		}
		
		time.Sleep(1 * time.Second)
	}
}

func sendSimpleTestMessage(message string, path []string) error {
	// 创建路径字符串
	pathString := strings.Join(path, ",")

	// 创建路由信息
	routeInfo := protocol.RouteInfo{
		TargetAddress: pathString,
		Priority:      1,
		TTL:           64,
		Flags:         protocol.RouteFlagDirect,
		Timestamp:     time.Now().Unix(),
	}

	// 创建数据块
	blockID := uint64(time.Now().UnixNano())
	dataBlock := protocol.CreateDataBlock(blockID, []byte(message), routeInfo)

	// 序列化数据块
	serialized, err := dataBlock.Serialize()
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}

	// 发送到第一跳
	firstHop := path[0]
	conn, err := net.DialTimeout("tcp", firstHop, 5*time.Second)
	if err != nil {
		return fmt.Errorf("连接第一跳失败: %v", err)
	}
	defer conn.Close()

	_, err = conn.Write(serialized)
	if err != nil {
		return fmt.Errorf("发送数据失败: %v", err)
	}

	fmt.Printf("📡 数据块已发送: BlockID=%d, 大小=%d字节\n", 
		blockID, len(serialized))

	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
