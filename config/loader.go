package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadConfig 从文件加载配置
func LoadConfig(path string) (*Config, error) {
	// 占位符实现，实际使用时需要完善
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file error: %v", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config error: %v", err)
	}

	// 设置默认值
	setDefaultValues(&cfg)

	return &cfg, nil
}

// setDefaultValues 设置配置的默认值
func setDefaultValues(cfg *Config) {
	// 设置通用配置默认值
	if cfg.Common.LogLevel == "" {
		cfg.Common.LogLevel = "info"
	}
	if cfg.Common.NodeID == "" {
		cfg.Common.NodeID = "proxy-node"
	}
	if cfg.Common.Region == "" {
		cfg.Common.Region = "default"
	}

	// 设置监听器默认值
	for i := range cfg.Listeners {
		if cfg.Listeners[i].Name == "" {
			cfg.Listeners[i].Name = fmt.Sprintf("listener-%d", i)
		}
	}

	// 设置调度器默认值
	if cfg.Scheduler.Address == "" {
		cfg.Scheduler.Address = "127.0.0.1:9999"
	}
	if cfg.Scheduler.ConnectTimeout == 0 {
		cfg.Scheduler.ConnectTimeout = 30 * time.Second
	}
	if cfg.Scheduler.RequestTimeout == 0 {
		cfg.Scheduler.RequestTimeout = 10 * time.Second
	}
	if cfg.Scheduler.RetryInterval == 0 {
		cfg.Scheduler.RetryInterval = 5 * time.Second
	}
	if cfg.Scheduler.MaxRetries == 0 {
		cfg.Scheduler.MaxRetries = 3
	}
	if cfg.Scheduler.HeartbeatInterval == 0 {
		cfg.Scheduler.HeartbeatInterval = 30 * time.Second
	}

	// 设置对等节点默认值
	for i := range cfg.Peers {
		if cfg.Peers[i].Type == "" {
			cfg.Peers[i].Type = "tcp"
		}
		if cfg.Peers[i].Region == "" {
			cfg.Peers[i].Region = "default"
		}
	}
}
