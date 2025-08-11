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

	fmt.Println("🚀 高性能数据面多跳转发演示 (并行TCP + smux多路复用)")
	fmt.Println(strings.Repeat("=", 70))

	// 启动多跳转发链路
	// Client → Proxy1(8430) → Proxy2(8431) → Proxy3(8432) → Target
	
	// 启动代理节点 (使用高性能连接)
	proxy1 := startHybridForwarderNode("127.0.0.1:8430", "8430", "HybridProxy1")
	proxy2 := startHybridForwarderNode("127.0.0.1:8431", "8431", "HybridProxy2") 
	proxy3 := startHybridForwarderNode("127.0.0.1:8432", "8432", "HybridProxy3")
	
	// 等待节点启动
	time.Sleep(1 * time.Second)
	
	// 演示高性能多跳转发
	demonstrateHybridMultiHopForwarding()
	
	// 显示统计信息
	time.Sleep(2 * time.Second)
	showHybridForwarderStats(proxy1, "HybridProxy1")
	showHybridForwarderStats(proxy2, "HybridProxy2") 
	showHybridForwarderStats(proxy3, "HybridProxy3")
	
	// 清理资源
	proxy1.Stop()
	proxy2.Stop()
	proxy3.Stop()
	
	fmt.Println("\n🎉 高性能数据面转发演示完成！")
}

func startHybridForwarderNode(nodeAddress, listenPort, nodeName string) *dataplane.DataPlaneForwarder {
	forwarder := dataplane.NewDataPlaneForwarder(nodeAddress, ":"+listenPort)
	
	if err := forwarder.Start(); err != nil {
		fmt.Printf("❌ %s 启动失败: %v\n", nodeName, err)
		return nil
	}
	
	fmt.Printf("🔗 %s 启动成功: %s (并行TCP + smux多路复用)\n", nodeName, nodeAddress)
	return forwarder
}

func demonstrateHybridMultiHopForwarding() {
	fmt.Println("\n🧪 测试高性能多跳数据转发")
	fmt.Println(strings.Repeat("-", 50))
	
	// 定义转发路径
	forwardingPath := []string{
		"127.0.0.1:8430", // HybridProxy1
		"127.0.0.1:8431", // HybridProxy2  
		"127.0.0.1:8432", // HybridProxy3 (最终目标)
	}
	
	fmt.Printf("📍 转发路径: %v\n", forwardingPath)
	fmt.Printf("🚀 使用技术: 并行TCP连接 + smux多路复用\n")
	
	// 测试数据 - 更多更大的数据来体现性能优势
	testMessages := []string{
		"High-performance message 1: Testing parallel TCP connections",
		"High-performance message 2: Testing smux multiplexing capabilities", 
		"High-performance message 3: Large data block to demonstrate throughput improvements with parallel connections and stream multiplexing",
		"High-performance message 4: Concurrent transmission test",
		"High-performance message 5: Final performance validation",
	}
	
	fmt.Printf("📤 准备发送 %d 条高性能测试消息\n", len(testMessages))
	
	// 并发发送多条消息来体现性能优势
	for i, message := range testMessages {
		fmt.Printf("\n📤 发送高性能消息 %d: %s\n", i+1, truncateMessage(message, 50))
		
		if err := sendHybridDataBlock(message, forwardingPath); err != nil {
			fmt.Printf("❌ 消息 %d 发送失败: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ 消息 %d 发送成功 (通过并行TCP + smux)\n", i+1)
		}
		
		// 短暂间隔以观察并发效果
		time.Sleep(100 * time.Millisecond)
	}
	
	// 测试大数据块传输
	fmt.Println("\n📦 测试大数据块传输 (1MB)")
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}
	
	if err := sendHybridDataBlock(string(largeData), forwardingPath); err != nil {
		fmt.Printf("❌ 大数据块发送失败: %v\n", err)
	} else {
		fmt.Printf("✅ 大数据块发送成功: %d字节 (通过并行TCP + smux)\n", len(largeData))
	}
}

func sendHybridDataBlock(message string, path []string) error {
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
	blockID := uint64(time.Now().UnixNano())
	dataBlock := protocol.CreateDataBlock(blockID, []byte(message), routeInfo)
	
	// 序列化数据块
	serialized, err := dataBlock.Serialize()
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}
	
	// 发送到第一跳 (Client端仍使用简单TCP连接)
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
	
	fmt.Printf("📡 数据块已发送: BlockID=%d, 大小=%d字节, 路径=%s\n", 
		blockID, len(serialized), pathString)
	
	return nil
}

func showHybridForwarderStats(forwarder *dataplane.DataPlaneForwarder, name string) {
	if forwarder == nil {
		return
	}
	
	stats := forwarder.GetStats()
	fmt.Printf("\n📊 %s 统计信息 (高性能模式):\n", name)
	fmt.Printf("  节点地址: %s\n", stats.NodeAddress)
	fmt.Printf("  总连接数: %d\n", stats.TotalConnections)
	fmt.Printf("  活跃连接: %d\n", stats.ActiveConnections)
	fmt.Printf("  转发次数: %d\n", stats.TotalForwards)
	fmt.Printf("  总字节数: %d\n", stats.TotalBytes)
	fmt.Printf("  错误次数: %d\n", stats.ErrorCount)
	fmt.Printf("  传输技术: 并行TCP + smux多路复用\n")
}

func truncateMessage(message string, maxLen int) string {
	if len(message) <= maxLen {
		return message
	}
	return message[:maxLen] + "..."
}
