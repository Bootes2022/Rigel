package config

import "time"

// Config 动态代理配置
type Config struct {
	Common    CommonConfig     `json:"common" yaml:"common"`
	Listeners []ListenerConfig `json:"listeners" yaml:"listeners"`
	Scheduler SchedulerConfig  `json:"scheduler" yaml:"scheduler"`
	Peers     []PeerConfig     `json:"peers" yaml:"peers"`
}

// CommonConfig 通用配置
type CommonConfig struct {
	LogLevel  string `json:"log_level" yaml:"log_level"`
	LogFile   string `json:"log_file" yaml:"log_file"`
	AdminAddr string `json:"admin_addr" yaml:"admin_addr"`
	NodeID    string `json:"node_id" yaml:"node_id"`
	Region    string `json:"region" yaml:"region"`
}

// ListenerConfig 监听器配置
type ListenerConfig struct {
	Name string `json:"name" yaml:"name"`
	Bind string `json:"bind" yaml:"bind"`
}

// SchedulerConfig 调度器配置
type SchedulerConfig struct {
	Address           string        `json:"address" yaml:"address"`
	ConnectTimeout    time.Duration `json:"connect_timeout" yaml:"connect_timeout"`
	RequestTimeout    time.Duration `json:"request_timeout" yaml:"request_timeout"`
	RetryInterval     time.Duration `json:"retry_interval" yaml:"retry_interval"`
	MaxRetries        int           `json:"max_retries" yaml:"max_retries"`
	HeartbeatInterval time.Duration `json:"heartbeat_interval" yaml:"heartbeat_interval"`
}

// PeerConfig 对等节点配置
type PeerConfig struct {
	Name    string `json:"name" yaml:"name"`
	Address string `json:"address" yaml:"address"`
	Type    string `json:"type,omitempty" yaml:"type,omitempty"` // tcp, tls等
	Region  string `json:"region,omitempty" yaml:"region,omitempty"`
}
