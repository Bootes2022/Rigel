# Rigel部署验证脚本

param(
    [string]$Password = "20010221cc."
)

Write-Host "🔍 Rigel部署验证脚本" -ForegroundColor Green
Write-Host "=====================" -ForegroundColor Green

# 服务器配置
$servers = @{
    "Proxy1" = @{
        "IP" = "47.101.37.95"
        "Description" = "日本东京"
    }
    "Proxy2" = @{
        "IP" = "101.201.209.244" 
        "Description" = "中国杭州"
    }
    "Target" = @{
        "IP" = "120.27.113.148"
        "Description" = "中国深圳"
    }
}

Write-Host "📋 验证服务器状态:" -ForegroundColor Yellow
Write-Host ""

foreach ($server in $servers.GetEnumerator()) {
    $serverName = $server.Key
    $serverIP = $server.Value.IP
    $description = $server.Value.Description
    
    Write-Host "🔍 检查 $serverName ($serverIP - $description)" -ForegroundColor Cyan
    
    # 1. 测试网络连通性
    try {
        $ping = Test-Connection -ComputerName $serverIP -Count 1 -Quiet
        if ($ping) {
            Write-Host "  ✅ 网络连通正常" -ForegroundColor Green
        } else {
            Write-Host "  ❌ 网络连通失败" -ForegroundColor Red
            continue
        }
    }
    catch {
        Write-Host "  ❌ 网络测试异常: $($_.Exception.Message)" -ForegroundColor Red
        continue
    }
    
    # 2. 测试端口连通性
    try {
        $connection = Test-NetConnection -ComputerName $serverIP -Port 8080 -WarningAction SilentlyContinue
        if ($connection.TcpTestSucceeded) {
            Write-Host "  ✅ 端口8080连通正常" -ForegroundColor Green
        } else {
            Write-Host "  ❌ 端口8080连接失败" -ForegroundColor Red
        }
    }
    catch {
        Write-Host "  ❌ 端口测试异常: $($_.Exception.Message)" -ForegroundColor Red
    }
    
    # 3. 检查进程状态
    try {
        $processCheck = ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 root@$serverIP "ps aux | grep rigel-proxy | grep -v grep | wc -l" 2>$null
        if ($processCheck -gt 0) {
            Write-Host "  ✅ rigel-proxy进程运行正常" -ForegroundColor Green
        } else {
            Write-Host "  ❌ rigel-proxy进程未运行" -ForegroundColor Red
        }
    }
    catch {
        Write-Host "  ❌ 进程检查失败" -ForegroundColor Red
    }
    
    # 4. 检查日志
    try {
        $logCheck = ssh -o StrictHostKeyChecking=no -o ConnectTimeout=5 root@$serverIP "tail -1 /opt/rigel/logs/rigel.log 2>/dev/null || echo 'No log'" 2>$null
        if ($logCheck -and $logCheck -ne "No log") {
            Write-Host "  ✅ 日志文件存在" -ForegroundColor Green
        } else {
            Write-Host "  ⚠️ 日志文件不存在或为空" -ForegroundColor Yellow
        }
    }
    catch {
        Write-Host "  ❌ 日志检查失败" -ForegroundColor Red
    }
    
    Write-Host ""
}

# 5. 运行连通性测试
Write-Host "🧪 运行连通性测试..." -ForegroundColor Magenta

try {
    Write-Host "执行: go run deploy\single-test.go" -ForegroundColor Gray
    $testResult = & go run deploy\single-test.go 2>&1
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "✅ 连通性测试成功" -ForegroundColor Green
        Write-Host $testResult -ForegroundColor Gray
    } else {
        Write-Host "❌ 连通性测试失败" -ForegroundColor Red
        Write-Host $testResult -ForegroundColor Red
    }
}
catch {
    Write-Host "❌ 测试执行异常: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎯 验证完成！" -ForegroundColor Green
