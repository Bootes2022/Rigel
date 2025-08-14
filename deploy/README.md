# Rigel TCP代理系统部署指南

## 📋 概述

Rigel是一个高性能的TCP代理系统，支持多跳数据转发。本文档描述了如何将系统部署到真实服务器环境。

## 🏗️ 系统架构

```
本地客户端
    ↓
Proxy1 (47.101.37.95:8080) - 日本东京
    ↓
Proxy2 (101.201.209.244:8080) - 中国杭州  
    ↓
Target (120.27.113.148:8080) - 中国深圳
```

## 📁 文件结构

```
deploy/
├── bin/                         # 编译后的二进制文件
│   ├── rigel-proxy.exe          # Windows版本
│   ├── rigel-proxy-linux        # Linux版本
│   └── rigel-client.exe         # 测试客户端
├── config/                      # 配置文件
│   ├── proxy1.yaml              # Proxy1配置 (日本东京)
│   ├── proxy2.yaml              # Proxy2配置 (中国杭州)
│   └── target.yaml              # Target配置 (中国深圳)
├── scripts/                     # 部署脚本
│   ├── deploy-fixed.ps1         # 主部署脚本
│   ├── verify-deployment.ps1    # 验证脚本
│   └── fix-target.ps1           # Target修复脚本
├── reports/                     # 测试报告
│   ├── test-report.md           # 详细技术报告
│   └── executive-summary.md     # 执行总结报告
├── real-world-client.go         # 真实世界测试客户端
├── single-test.go               # 单次测试客户端
├── deployment-checklist.md     # 部署检查清单
└── README.md                    # 本文档
```

## 🚀 部署步骤

### 前置条件

1. **本地环境**：
   - Windows 10/11
   - PowerShell 5.0+
   - Go 1.19+
   - SSH客户端

2. **服务器环境**：
   - Ubuntu 20.04+ 
   - Root访问权限
   - 开放8080端口

3. **网络要求**：
   - 服务器间网络连通
   - 防火墙允许8080端口

### 第1步：编译程序

```bash
# 编译Windows版本
go build -o deploy/bin/rigel-proxy.exe cmd/proxy/main.go

# 编译Linux版本  
set GOOS=linux
set GOARCH=amd64
go build -o deploy/bin/rigel-proxy-linux cmd/proxy/main.go

# 编译测试客户端
go build -o deploy/bin/rigel-client.exe deploy/real-world-client.go
```

### 第2步：配置文件

配置文件已预先创建在 `config/` 目录下：

- `proxy1.yaml` - Proxy1配置
- `proxy2.yaml` - Proxy2配置  
- `target.yaml` - Target配置

### 第3步：执行部署

```powershell
# 主部署脚本
.\deploy\scripts\deploy-fixed.ps1

# 如果Target有问题，运行修复脚本
.\deploy\scripts\fix-target.ps1

# 验证部署状态
.\deploy\scripts\verify-deployment.ps1
```

### 第4步：测试系统

```bash
# 单次测试
go run deploy/single-test.go

# 完整测试
go run deploy/real-world-client.go
```

## 🔧 故障排除

### 常见问题

1. **SSH连接失败**
   - 检查服务器IP和密码
   - 确认SSH服务运行正常
   - 检查网络连通性

2. **端口连接失败**
   - 检查防火墙设置
   - 确认程序正确启动
   - 验证配置文件正确

3. **程序启动失败**
   - 检查程序权限
   - 验证配置文件语法
   - 查看日志文件

### 调试命令

```bash
# 检查进程状态
ssh root@<server_ip> "ps aux | grep rigel-proxy"

# 检查端口监听
ssh root@<server_ip> "netstat -tlnp | grep 8080"

# 查看日志
ssh root@<server_ip> "tail -f /opt/rigel/logs/rigel.log"

# 手动启动程序
ssh root@<server_ip> "cd /opt/rigel && ./bin/rigel-proxy -c config.yaml"
```

## 📊 性能监控

### 关键指标

- **连接延迟**：通常 < 1ms
- **数据传输速度**：取决于网络带宽
- **内存使用**：通常 < 50MB
- **CPU使用**：通常 < 5%

### 监控命令

```bash
# 系统资源使用
ssh root@<server_ip> "top -p \$(pgrep rigel-proxy)"

# 网络连接状态
ssh root@<server_ip> "ss -tuln | grep 8080"

# 磁盘使用
ssh root@<server_ip> "df -h /opt/rigel"
```

## 🔒 安全考虑

1. **防火墙配置**：只开放必要端口
2. **访问控制**：限制SSH访问来源
3. **日志监控**：定期检查访问日志
4. **更新维护**：定期更新系统和程序

## 📞 支持

如遇问题，请检查：

1. 部署日志输出
2. 服务器系统日志
3. 程序运行日志
4. 网络连通性测试

## 📝 更新日志

- **v2.0** - 修复部署脚本，增加错误处理
- **v1.0** - 初始版本，基本部署功能
