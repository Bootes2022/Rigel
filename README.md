# 动态 TCP 代理 (Dynamic TCP Proxy)

一个基于 Go 语言开发的动态 TCP 代理服务器，实现了完全基于调度指令的路由决策，支持与 Traffic Scheduler 的实时通信和智能路径选择。

## 🚀 重大架构升级

本项目已完成从静态路由到动态调度的重大架构重构，现在完全符合 PRD 中轻量级代理模块的要求。

## 核心特性

### 📡 **动态调度路由**
- **调度指令驱动**：完全基于 Traffic Scheduler 的调度指令进行路由决策
- **实时路径选择**：根据网络状况和优化目标动态选择最优路径
- **智能缓存**：本地缓存调度指令，提高路由决策性能
- **多路径支持**：支持单个流的多路径传输和负载均衡

### 🔄 **统一代理模型**
- **对等节点**：每个代理节点功能相同，无客户端/服务端区分
- **双向转发**：既可接收连接，也可主动发起连接
- **透明中继**：作为中间节点转发数据

### 🌐 **调度器集成**
- **实时通信**：与 Traffic Scheduler 保持持续连接
- **状态报告**：定期报告流量统计和节点状态
- **指令订阅**：订阅调度更新，支持路径动态调整
- **故障恢复**：调度器不可用时的优雅降级

## 架构设计

```
Traffic Scheduler (控制平面)
        ↓ 调度指令
应用A ←→ 代理1 ←→ 代理2 ←→ 应用B
        ↑ 状态报告
```

### 核心组件

- **动态路由引擎 (DynamicRouteEngine)**：基于调度指令的智能路由决策
- **调度器客户端 (SchedulerClient)**：与 Traffic Scheduler 的通信接口
- **流管理器 (FlowManager)**：数据流生命周期管理和统计
- **指令缓存 (InstructionCache)**：调度指令的本地缓存和更新
- **对等节点管理器 (PeerManager)**：管理对等节点连接和健康状态

## 配置格式

### 动态代理配置示例

```yaml
# 通用配置
common:
  log_level: "info"
  log_file: "dynamic-proxy.log"
  admin_addr: "127.0.0.1:8888"
  node_id: "proxy-node-1"
  region: "us-east"

# 监听器配置
listeners:
  - name: "web-listener"
    bind: "0.0.0.0:8080"
  - name: "api-listener"
    bind: "0.0.0.0:9090"

# 调度器配置
scheduler:
  address: "127.0.0.1:9999"
  connect_timeout: "30s"
  request_timeout: "10s"
  retry_interval: "5s"
  max_retries: 3
  heartbeat_interval: "30s"

# 对等节点配置（用于节点发现和通信）
peers:
  - name: "proxy-node-2"
    address: "192.168.1.50:9090"
    type: "tcp"
    region: "us-west"

  - name: "proxy-node-3"
    address: "192.168.1.51:9090"
    type: "tcp"
    region: "eu-central"
```

### 配置说明

#### 通用配置 (common)
- `node_id`: 节点唯一标识符
- `region`: 节点所在区域
- `log_level`: 日志级别
- `log_file`: 日志文件路径

#### 监听器配置 (listeners)
- `name`: 监听器名称
- `bind`: 监听地址和端口

#### 调度器配置 (scheduler)
- `address`: Traffic Scheduler 地址
- `connect_timeout`: 连接超时时间
- `request_timeout`: 请求超时时间
- `retry_interval`: 重试间隔
- `max_retries`: 最大重试次数
- `heartbeat_interval`: 心跳间隔

#### 对等节点配置 (peers)
- `name`: 节点名称
- `address`: 节点地址
- `type`: 连接类型（默认 tcp）
- `region`: 节点区域

## 使用方法

### 编译

```bash
go build ./cmd/proxy
```

### 运行

```bash
# 使用默认配置文件
./proxy

# 指定配置文件
./proxy -c config.yaml
```

### 测试

```bash
# 运行单元测试
go test ./internal/router -v
go test ./internal/peer -v

# 运行集成测试
go test ./test -v
```

## 使用场景

### 1. 负载均衡
```yaml
listeners:
  - name: "lb-listener"
    bind: "0.0.0.0:80"

routes:
  - name: "web-route"
    match:
      port: 80
    target: "backend-server:8080"
```

### 2. 代理链
```yaml
# 节点1配置
routes:
  - name: "chain-route"
    match:
      port: 8080
    target: "proxy-node-2"

peers:
  - name: "proxy-node-2"
    address: "192.168.1.50:9090"
```

### 3. 服务网格
```yaml
# 多个监听器和路由
listeners:
  - name: "service-a"
    bind: "0.0.0.0:8080"
  - name: "service-b"
    bind: "0.0.0.0:9090"

routes:
  - name: "route-a"
    match:
      port: 8080
    target: "service-a-backend:8080"
  - name: "route-b"
    match:
      port: 9090
    target: "service-b-backend:9090"
```

## API 接口

### 统计信息
```go
stats := proxy.GetStats()
// 返回：连接数、流量统计、节点信息等
```

### 监听器状态
```go
listeners := proxy.GetListeners()
// 返回：监听器列表和状态
```

### 对等节点状态
```go
peers := proxy.GetPeers()
// 返回：对等节点列表和健康状态
```

## 性能特性

- ✅ 支持数百个并发连接
- ✅ 低延迟数据转发
- ✅ 高效的内存和 CPU 使用
- ✅ 优雅的连接关闭
- ✅ 完整的错误处理

## 监控和日志

### 日志级别
- `debug`: 详细调试信息
- `info`: 一般信息
- `warn`: 警告信息
- `error`: 错误信息
- `fatal`: 致命错误

### 统计指标
- 当前连接数
- 总连接数
- 入站/出站流量
- 节点健康状态
- 路由匹配统计

## 开发和扩展

### 项目结构
```
├── cmd/proxy/          # 主程序入口
├── config/             # 配置管理
├── internal/
│   ├── proxy/          # 代理核心
│   ├── router/         # 路由引擎
│   ├── peer/           # 对等节点管理
│   ├── listener/       # 监听器抽象
│   ├── transport/      # 传输层抽象
│   └── metrics/        # 指标收集
├── pkg/
│   ├── log/            # 日志系统
│   └── util/           # 工具函数
└── test/               # 测试用例
```

### 扩展点
- 支持更多传输协议（TLS、UDP等）
- 添加认证和授权机制
- 实现配置热重载
- 添加 Web 管理界面
- 支持插件系统

## 许可证

MIT License
