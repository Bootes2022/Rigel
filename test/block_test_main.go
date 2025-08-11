package main

import (
	"fmt"
	"tcp-proxy/internal/protocol"
)

func main() {
	fmt.Println("🧪 最小化块协议测试")
	
	// 创建客户端
	client := protocol.NewBlockProtocolClient()
	fmt.Println("✅ 客户端创建成功")
	
	// 编码数据
	testData := []byte("Hello World")
	block, err := client.EncodeDataWithRoute(testData, "127.0.0.1:8080", 1)
	if err != nil {
		fmt.Printf("❌ 编码失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 数据编码成功: BlockID=%d\n", block.Header.BlockID)
	
	// 序列化
	serialized, err := block.Serialize()
	if err != nil {
		fmt.Printf("❌ 序列化失败: %v\n", err)
		return
	}
	fmt.Printf("✅ 序列化成功: %d 字节\n", len(serialized))
	
	// 创建服务端
	server := protocol.NewBlockProtocolServer()
	fmt.Println("✅ 服务端创建成功")
	
	// 处理数据
	blocks, err := server.ProcessIncomingData(serialized)
	if err != nil {
		fmt.Printf("❌ 处理失败: %v\n", err)
		return
	}
	
	if len(blocks) == 1 {
		fmt.Printf("✅ 解析成功: %s -> %s\n", 
			string(blocks[0].Data), 
			blocks[0].Header.RouteInfo.TargetAddress)
	} else {
		fmt.Printf("❌ 期望1个块，得到%d个\n", len(blocks))
	}
	
	fmt.Println("🎉 测试完成！")
}
