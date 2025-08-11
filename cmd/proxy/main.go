package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"tcp-proxy/config"
	"tcp-proxy/internal/dataplane"
	"tcp-proxy/pkg/log"
)

var (
	configFile = flag.String("c", "config.yaml", "config file path")
)

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := log.Init(cfg.Common.LogLevel, cfg.Common.LogFile); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	log.Infof("Starting data plane forwarder node: %s", cfg.Common.NodeID)

	// 创建数据面转发器
	nodeAddress := cfg.Common.NodeID + ":8080" // 默认端口
	forwarder := dataplane.NewDataPlaneForwarder(nodeAddress, ":8080")

	// 启动转发器
	if err := forwarder.Start(); err != nil {
		log.Fatalf("Failed to start data plane forwarder: %v", err)
	}

	// 等待信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	// 停止转发器
	log.Info("Shutting down...")
	if err := forwarder.Stop(); err != nil {
		log.Errorf("Error stopping data plane forwarder: %v", err)
	}
	log.Info("Data plane forwarder stopped")
}
