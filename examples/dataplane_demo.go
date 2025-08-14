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

	fmt.Println("🚀 数据面多跳转发演示")
	fmt.Println(strings.Repeat("=", 60))

	// 启动多跳转发链路
	// Client → Proxy1(8420) → Proxy2(8421) → Proxy3(8422) → Target
	
	// 启动代理节点
	proxy1 := startForwarderNode("127.0.0.1:8420", "8420", "proxy1")
	proxy2 := startForwarderNode("127.0.0.1:8421", "8421", "proxy2") 
	proxy3 := startForwarderNode("127.0.0.1:8422", "8422", "proxy3")
	
	// 等待节点启动
	time.Sleep(500 * time.Millisecond)
	
	// 演示多跳转发
	demonstrateMultiHopForwarding()
	
	// 显示统计信息
	time.Sleep(1 * time.Second)
	showForwarderStats(proxy1, "Proxy1")
	showForwarderStats(proxy2, "Proxy2") 
	showForwarderStats(proxy3, "Proxy3")
	
	// 清理资源
	proxy1.Stop()
	proxy2.Stop()
	proxy3.Stop()
	
	fmt.Println("\n🎉 数据面转发演示完成！")
}

func startForwarderNode(nodeAddress, listenPort, nodeName string) *dataplane.DataPlaneForwarder {
	forwarder := dataplane.NewDataPlaneForwarder(nodeAddress, ":"+listenPort)
	
	if err := forwarder.Start(); err != nil {
		fmt.Printf("❌ %s 启动失败: %v\n", nodeName, err)
		return nil
	}
	
	fmt.Printf("🔗 %s 启动成功: %s\n", nodeName, nodeAddress)
	return forwarder
}

func demonstrateMultiHopForwarding() {
	fmt.Println("\n🧪 测试多跳数据转发")
	fmt.Println(strings.Repeat("-", 40))
	
	// 定义转发路径
	forwardingPath := []string{
		"127.0.0.1:8420", // Proxy1
		"127.0.0.1:8421", // Proxy2  
		"127.0.0.1:8422", // Proxy3 (最终目标)
	}
	
	fmt.Printf("📍 转发路径: %v\n", forwardingPath)
	
	// 测试数据
	testMessages := []string{
		"Hello Multi-Hop Forwarding!",
		"This is message number 2",
		"Final test message with some longer content to verify data integrity",
	}
	
	for i, message := range testMessages {
		fmt.Printf("\n📤 发送消息 %d: %s\n", i+1, message)
		
		if err := sendDataBlock(message, forwardingPath); err != nil {
			fmt.Printf("❌ 消息 %d 发送失败: %v\n", i+1, err)
		} else {
			fmt.Printf("✅ 消息 %d 发送成功\n", i+1)
		}
		
		time.Sleep(200 * time.Millisecond)
	}
}

func sendDataBlock(message string, path []string) error {
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
	
	fmt.Printf("📡 数据块已发送: BlockID=%d, 大小=%d字节, 路径=%s\n", 
		blockID, len(serialized), pathString)
	
	return nil
}

func showForwarderStats(forwarder *dataplane.DataPlaneForwarder, name string) {
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
}
