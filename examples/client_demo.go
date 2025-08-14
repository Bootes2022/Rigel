package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"tcp-proxy/internal/client"
	"tcp-proxy/pkg/log"
)

func main() {
	// 初始化日志
	log.Init("info", "")

	fmt.Println("🚀 Rigel Client模块演示")
	fmt.Println(strings.Repeat("=", 60))

	// 启动模拟的代理节点（提供路径计算服务）
	go startMockProxyNode("127.0.0.1:8400")
	go startMockProxyNode("127.0.0.1:8401")
	go startMockProxyNode("127.0.0.1:8402")
	time.Sleep(200 * time.Millisecond)

	// 启动目标服务器
	go startTargetServer("127.0.0.1:8500")
	time.Sleep(100 * time.Millisecond)

	// 演示Client模块功能
	demonstrateClientModule()

	fmt.Println("\n🎉 Client模块演示完成！")
}

func demonstrateClientModule() {
	fmt.Println("\n🧪 测试Rigel Client模块")
	fmt.Println(strings.Repeat("-", 40))

	// 创建客户端配置
	config := client.DefaultClientConfig()
	config.ClientID = "demo_client_001"
	config.ChunkSize = 32 * 1024 // 32KB chunks
	config.MaxConcurrency = 5

	// 创建Rigel客户端
	rigelClient, err := client.NewRigelClient(config)
	if err != nil {
		fmt.Printf("❌ 创建客户端失败: %v\n", err)
		return
	}
	defer rigelClient.Close()

	fmt.Printf("✅ Rigel客户端创建成功: %s\n", config.ClientID)

	// 测试1: 发送小数据
	fmt.Println("\n1️⃣ 发送小数据测试")
	smallData := "Hello Rigel! This is a small message for testing."
	if err := rigelClient.SendData([]byte(smallData), "127.0.0.1:8500"); err != nil {
		fmt.Printf("❌ 发送小数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 小数据发送成功: %d字节\n", len(smallData))
	}

	time.Sleep(500 * time.Millisecond)

	// 测试2: 发送大数据（会被拆分为多个块）
	fmt.Println("\n2️⃣ 发送大数据测试")
	largeData := make([]byte, 100*1024) // 100KB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	if err := rigelClient.SendDataWithPriority(largeData, "127.0.0.1:8500", 2); err != nil {
		fmt.Printf("❌ 发送大数据失败: %v\n", err)
	} else {
		fmt.Printf("✅ 大数据发送成功: %d字节\n", len(largeData))
	}

	time.Sleep(500 * time.Millisecond)

	// 测试3: 批量发送
	fmt.Println("\n3️⃣ 批量发送测试")
	for i := 0; i < 3; i++ {
		message := fmt.Sprintf("Batch message %d - timestamp: %d", i, time.Now().UnixNano())
		priority := uint8(i + 1)

		if err := rigelClient.SendDataWithPriority([]byte(message), "127.0.0.1:8500", priority); err != nil {
			fmt.Printf("❌ 批量消息 %d 发送失败: %v\n", i, err)
		} else {
			fmt.Printf("✅ 批量消息 %d 发送成功 (优先级=%d)\n", i, priority)
		}

		time.Sleep(200 * time.Millisecond)
	}

	// 显示统计信息
	time.Sleep(1 * time.Second)
	stats := rigelClient.GetStats()
	fmt.Printf("\n📊 客户端统计信息:\n")
	fmt.Printf("  客户端ID: %s\n", stats.ClientID)
	fmt.Printf("  任务统计:\n")
	fmt.Printf("    总任务数: %d\n", stats.TaskStats.TotalTasks)
	fmt.Printf("    总数据块: %d\n", stats.TaskStats.TotalChunks)
	fmt.Printf("    总字节数: %d\n", stats.TaskStats.TotalBytes)
	fmt.Printf("    平均块数: %.2f\n", stats.TaskStats.AverageChunks)
	fmt.Printf("  决策统计:\n")
	fmt.Printf("    总请求数: %d\n", stats.DecisionStats.TotalRequests)
	fmt.Printf("    成功率: %.2f%%\n", stats.DecisionStats.SuccessRate*100)
	fmt.Printf("    平均延迟: %.2fms\n", stats.DecisionStats.AverageLatency)
	fmt.Printf("  速率限制:\n")
	fmt.Printf("    总请求数: %d\n", stats.RateLimiterStats.TotalRequests)
	fmt.Printf("    通过率: %.2f%%\n", stats.RateLimiterStats.GrantRate*100)
	fmt.Printf("    当前速率: %.2f tokens/sec\n", stats.RateLimiterStats.CurrentRate)
}

// startMockProxyNode 启动模拟代理节点（提供路径计算服务）
func startMockProxyNode(address string) {
	mux := http.NewServeMux()

	// 路径计算API
	mux.HandleFunc("/api/v1/path/compute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// 模拟路径计算（简化实现）

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// 简单的JSON响应
		fmt.Fprintf(w, `{
			"path": ["%s", "127.0.0.1:8500"],
			"recommended_rate": 1048576,
			"ttl": 300,
			"request_id": "req_%d"
		}`, address, time.Now().UnixNano())

		fmt.Printf("📡 代理节点 %s 提供路径计算服务\n", address)
	})

	fmt.Printf("🖥️  模拟代理节点启动: %s\n", address)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("❌ 代理节点 %s 启动失败: %v\n", address, err)
	}
}

// startTargetServer 启动目标服务器
func startTargetServer(address string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("📨 目标服务器接收请求: %s %s\n", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Response from target server %s", address)
	})

	fmt.Printf("🎯 目标服务器启动: %s\n", address)

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Printf("❌ 目标服务器 %s 启动失败: %v\n", address, err)
	}
}
