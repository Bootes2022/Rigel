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

	fmt.Println("🌐 真实世界3跳转发测试")
	fmt.Println("=================================")
	fmt.Println("本地电脑 → Proxy1 → Proxy2 → Target")
	fmt.Println()

	// 真实服务器路径 - 完整3跳转发
	realWorldPath := []string{
		"47.101.37.95:8080",    // Proxy1
		"101.201.209.244:8080", // Proxy2
		"120.27.113.148:8080",  // Target
	}

	fmt.Printf("📍 真实转发路径:\n")
	fmt.Printf("  1. Proxy1: %s\n", realWorldPath[0])
	fmt.Printf("  2. Proxy2: %s\n", realWorldPath[1])
	fmt.Printf("  3. Target: %s\n", realWorldPath[2])
	fmt.Println()

	// 执行真实世界测试
	runRealWorldTests(realWorldPath)

	fmt.Println("\n🎉 真实世界测试完成！")
}

func runRealWorldTests(path []string) {
	fmt.Println("🧪 开始真实世界传输测试")
	fmt.Println(strings.Repeat("-", 40))

	// 测试1: 简单消息
	fmt.Println("📤 测试1: 简单消息传输")
	message1 := "Hello from local computer! Testing real-world 3-hop forwarding across Tokyo → Singapore → Hong Kong"
	if err := sendRealWorldMessage(message1, path, 1); err != nil {
		fmt.Printf("❌ 测试1失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试1成功: 简单消息传输\n")
	}
	time.Sleep(2 * time.Second)

	// 测试2: 带时间戳的消息
	fmt.Println("\n📤 测试2: 带时间戳的消息")
	message2 := fmt.Sprintf("Real-world test message sent at %s from local computer", time.Now().Format("2006-01-02 15:04:05"))
	if err := sendRealWorldMessage(message2, path, 2); err != nil {
		fmt.Printf("❌ 测试2失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试2成功: 带时间戳消息传输\n")
	}
	time.Sleep(2 * time.Second)

	// 测试3: 较大数据
	fmt.Println("\n📤 测试3: 较大数据传输")
	largeMessage := strings.Repeat("This is a larger test message for real-world 3-hop forwarding validation. ", 100)
	if err := sendRealWorldMessage(largeMessage, path, 3); err != nil {
		fmt.Printf("❌ 测试3失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试3成功: 较大数据传输 (%d字节)\n", len(largeMessage))
	}
	time.Sleep(2 * time.Second)

	// 测试4: 连续传输
	fmt.Println("\n📤 测试4: 连续传输测试")
	for i := 1; i <= 3; i++ {
		message := fmt.Sprintf("Continuous test message #%d - Real-world multi-hop forwarding", i)
		if err := sendRealWorldMessage(message, path, 10+i); err != nil {
			fmt.Printf("❌ 连续测试%d失败: %v\n", i, err)
		} else {
			fmt.Printf("✅ 连续测试%d成功\n", i)
		}
		time.Sleep(1 * time.Second)
	}
}

func sendRealWorldMessage(message string, path []string, testID int) error {
	// 创建路径字符串
	pathString := strings.Join(path, ",")

	// 创建路由信息
	routeInfo := protocol.RouteInfo{
		TargetAddress: pathString,
		Priority:      1,
		TTL:           64,
		Flags:         protocol.RouteFlagDirect | protocol.RouteFlagMultiPath,
		Timestamp:     time.Now().Unix(),
	}

	// 创建数据块
	blockID := uint64(time.Now().UnixNano() + int64(testID))
	dataBlock := protocol.CreateDataBlock(blockID, []byte(message), routeInfo)

	// 序列化数据块
	serialized, err := dataBlock.Serialize()
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}

	// 连接到第一跳 (日本东京)
	firstHop := path[0]
	fmt.Printf("🔗 连接到日本东京服务器: %s\n", firstHop)

	conn, err := net.DialTimeout("tcp", firstHop, 10*time.Second)
	if err != nil {
		return fmt.Errorf("连接日本东京服务器失败: %v", err)
	}
	defer conn.Close()

	// 发送数据
	startTime := time.Now()
	_, err = conn.Write(serialized)
	if err != nil {
		return fmt.Errorf("发送数据失败: %v", err)
	}
	duration := time.Since(startTime)

	fmt.Printf("📡 测试%d数据已发送: BlockID=%d, 大小=%d字节, 耗时=%v\n",
		testID, blockID, len(serialized), duration)

	// 等待一小段时间确保传输完成
	time.Sleep(500 * time.Millisecond)

	return nil
}

// 网络连通性测试
func testConnectivity(servers []string) {
	fmt.Println("🌐 测试网络连通性")

	for i, server := range servers {
		fmt.Printf("测试连接到服务器%d (%s)...", i+1, server)

		conn, err := net.DialTimeout("tcp", server, 5*time.Second)
		if err != nil {
			fmt.Printf(" ❌ 连接失败: %v\n", err)
		} else {
			conn.Close()
			fmt.Printf(" ✅ 连接成功\n")
		}
	}
}
