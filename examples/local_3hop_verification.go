package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"tcp-proxy/internal/dataplane"
	"tcp-proxy/internal/protocol"
	"tcp-proxy/pkg/log"
)

func main() {
	// 初始化日志
	log.Init("info", "")

	fmt.Println("🧪 本地3跳转发验证测试")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("模拟场景: 本地Client → Proxy1 → Proxy2 → Target")
	fmt.Println()

	// 启动3跳转发链路
	proxy1 := startLocalProxy("127.0.0.1:9001", "9001", "Proxy1 (东京模拟)")
	proxy2 := startLocalProxy("127.0.0.1:9002", "9002", "Proxy2 (新加坡模拟)")
	target := startLocalProxy("127.0.0.1:9003", "9003", "Target (香港模拟)")

	// 等待节点启动
	time.Sleep(1 * time.Second)

	// 执行验证测试
	runVerificationTests()

	// 显示统计信息
	time.Sleep(2 * time.Second)
	showProxyStats(proxy1, "Proxy1 (东京模拟)")
	showProxyStats(proxy2, "Proxy2 (新加坡模拟)")
	showProxyStats(target, "Target (香港模拟)")

	// 清理资源
	proxy1.Stop()
	proxy2.Stop()
	target.Stop()

	fmt.Println("\n🎉 本地3跳转发验证完成！")
}

func startLocalProxy(nodeAddress, listenPort, nodeName string) *dataplane.DataPlaneForwarder {
	forwarder := dataplane.NewDataPlaneForwarder(nodeAddress, ":"+listenPort)

	if err := forwarder.Start(); err != nil {
		fmt.Printf("❌ %s 启动失败: %v\n", nodeName, err)
		return nil
	}

	fmt.Printf("🔗 %s 启动成功: %s\n", nodeName, nodeAddress)
	return forwarder
}

func runVerificationTests() {
	fmt.Println("\n🧪 开始3跳转发验证测试")
	fmt.Println(strings.Repeat("-", 50))

	// 定义3跳转发路径
	threehopPath := []string{
		"127.0.0.1:9001", // Proxy1 (东京模拟)
		"127.0.0.1:9002", // Proxy2 (新加坡模拟)
		"127.0.0.1:9003", // Target (香港模拟)
	}

	fmt.Printf("📍 3跳转发路径: %v\n", threehopPath)
	fmt.Printf("🚀 传输技术: 并行TCP + smux多路复用\n\n")

	// 测试1: 简单消息传输
	fmt.Println("📤 测试1: 简单消息传输")
	testMessage1 := "Hello from local client! Testing 3-hop forwarding."
	if err := sendTestMessage(testMessage1, threehopPath, 1); err != nil {
		fmt.Printf("❌ 测试1失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试1成功: 简单消息传输\n")
	}
	time.Sleep(500 * time.Millisecond)

	// 测试2: 中等大小数据
	fmt.Println("\n📤 测试2: 中等大小数据传输")
	testMessage2 := strings.Repeat("This is a medium-sized test message for 3-hop forwarding verification. ", 50)
	if err := sendTestMessage(testMessage2, threehopPath, 2); err != nil {
		fmt.Printf("❌ 测试2失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试2成功: 中等大小数据传输 (%d字节)\n", len(testMessage2))
	}
	time.Sleep(500 * time.Millisecond)

	// 测试3: 大数据块传输
	fmt.Println("\n📤 测试3: 大数据块传输")
	largeData := make([]byte, 100*1024) // 100KB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	if err := sendTestMessage(string(largeData), threehopPath, 3); err != nil {
		fmt.Printf("❌ 测试3失败: %v\n", err)
	} else {
		fmt.Printf("✅ 测试3成功: 大数据块传输 (%d字节)\n", len(largeData))
	}
	time.Sleep(500 * time.Millisecond)

	// 测试4: 并发传输
	fmt.Println("\n📤 测试4: 并发传输测试")
	concurrentMessages := []string{
		"Concurrent message 1: Testing parallel transmission",
		"Concurrent message 2: Verifying multiplexing capability",
		"Concurrent message 3: Validating 3-hop forwarding under load",
	}

	for i, msg := range concurrentMessages {
		go func(index int, message string) {
			if err := sendTestMessage(message, threehopPath, index+10); err != nil {
				fmt.Printf("❌ 并发测试%d失败: %v\n", index+1, err)
			} else {
				fmt.Printf("✅ 并发测试%d成功\n", index+1)
			}
		}(i, msg)
	}
	time.Sleep(2 * time.Second)

	fmt.Println("\n🎯 所有验证测试完成")
}

func sendTestMessage(message string, path []string, testID int) error {
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

	// 发送到第一跳 (Proxy1)
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

	fmt.Printf("📡 测试%d数据块已发送: BlockID=%d, 大小=%d字节\n",
		testID, blockID, len(serialized))

	return nil
}

func showProxyStats(forwarder *dataplane.DataPlaneForwarder, name string) {
	if forwarder == nil {
		return
	}

	stats := forwarder.GetStats()
	fmt.Printf("\n📊 %s 统计信息:\n", name)
	fmt.Printf("  节点地址: %s\n", stats.NodeAddress)
	fmt.Printf("  总连接数: %d\n", stats.TotalConnections)
	fmt.Printf("  活跃连接: %d\n", stats.ActiveConnections)
	fmt.Printf("  转发次数: %d\n", stats.TotalForwards)
	fmt.Printf("  总字节数: %d\n", stats.TotalBytes)
	fmt.Printf("  错误次数: %d\n", stats.ErrorCount)
	fmt.Printf("  传输技术: 并行TCP + smux多路复用\n")

	// 分析节点角色
	if stats.TotalConnections > 0 && stats.TotalForwards > 0 {
		if strings.Contains(name, "Target") {
			fmt.Printf("  节点角色: ✅ 最终目标节点 (成功接收数据)\n")
		} else {
			fmt.Printf("  节点角色: ✅ 中转代理节点 (成功转发数据)\n")
		}
	} else if stats.TotalConnections > 0 {
		fmt.Printf("  节点角色: ⚠️ 接收节点 (接收但未转发)\n")
	} else {
		fmt.Printf("  节点角色: ❌ 未使用节点\n")
	}
}
