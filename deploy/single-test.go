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

	fmt.Println("🧪 单次2跳转发测试")
	fmt.Println("===================")

	// 真实服务器路径
	realPath := []string{
		"47.101.37.95:8080",    // Proxy1
		"101.201.209.244:8080", // Proxy2
	}

	fmt.Printf("📍 转发路径: %s → %s\n", realPath[0], realPath[1])
	fmt.Println()

	// 测试消息
	message := "Hello from local computer! Testing single 2-hop forwarding"
	
	fmt.Printf("📤 发送测试消息: %s\n", message)
	
	if err := sendMessage(message, realPath); err != nil {
		fmt.Printf("❌ 测试失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试成功!\n")
	}
}

func sendMessage(message string, path []string) error {
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
	fmt.Printf("🔗 连接到: %s\n", firstHop)
	
	conn, err := net.DialTimeout("tcp", firstHop, 10*time.Second)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer conn.Close()

	startTime := time.Now()
	_, err = conn.Write(serialized)
	if err != nil {
		return fmt.Errorf("发送数据失败: %v", err)
	}
	duration := time.Since(startTime)

	fmt.Printf("📡 数据已发送: BlockID=%d, 大小=%d字节, 耗时=%v\n",
		blockID, len(serialized), duration)

	return nil
}
