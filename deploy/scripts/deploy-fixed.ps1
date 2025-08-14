# Rigel固定版部署脚本

param(
    [string]$Password = "20010221cc."
)

Write-Host "🌐 Rigel TCP代理系统部署脚本 v2.0" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green

# 服务器配置
$servers = @{
    "Proxy1" = @{
        "IP" = "47.101.37.95"
        "Config" = "proxy1.yaml"
        "Description" = "日本东京"
    }
    "Proxy2" = @{
        "IP" = "101.201.209.244" 
        "Config" = "proxy2.yaml"
        "Description" = "中国杭州"
    }
    "Target" = @{
        "IP" = "120.27.113.148"
        "Config" = "target.yaml"
        "Description" = "中国深圳"
    }
}

Write-Host "📋 服务器列表:" -ForegroundColor Yellow
foreach ($server in $servers.GetEnumerator()) {
    Write-Host "  $($server.Key): $($server.Value.IP) ($($server.Value.Description))" -ForegroundColor Cyan
}
Write-Host ""

# 函数：执行SSH命令
function Invoke-SSHCommandSafe {
    param(
        [string]$Server,
        [string]$Command,
        [string]$Description
    )
    
    Write-Host "🔧 $Description..." -ForegroundColor Yellow
    
    try {
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = "ssh"
        $psi.Arguments = "-o StrictHostKeyChecking=no -o ConnectTimeout=10 root@$Server `"$Command`""
        $psi.RedirectStandardOutput = $true
        $psi.RedirectStandardError = $true
        $psi.UseShellExecute = $false
        $psi.CreateNoWindow = $true
        
        $process = New-Object System.Diagnostics.Process
        $process.StartInfo = $psi
        $process.Start() | Out-Null
        
        # 发送密码
        Start-Sleep -Milliseconds 500
        $process.StandardInput.WriteLine($Password)
        
        $process.WaitForExit(30000) # 30秒超时
        
        if ($process.ExitCode -eq 0) {
            Write-Host "✅ 成功: $Description" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ 失败: $Description" -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Host "❌ 异常: $Description - $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# 函数：上传文件
function Copy-FileToServerSafe {
    param(
        [string]$Server,
        [string]$LocalPath,
        [string]$RemotePath,
        [string]$Description
    )
    
    Write-Host "📤 $Description..." -ForegroundColor Yellow
    
    try {
        $psi = New-Object System.Diagnostics.ProcessStartInfo
        $psi.FileName = "scp"
        $psi.Arguments = "-o StrictHostKeyChecking=no `"$LocalPath`" root@${Server}:$RemotePath"
        $psi.RedirectStandardOutput = $true
        $psi.RedirectStandardError = $true
        $psi.UseShellExecute = $false
        $psi.CreateNoWindow = $true
        
        $process = New-Object System.Diagnostics.Process
        $process.StartInfo = $psi
        $process.Start() | Out-Null
        
        # 发送密码
        Start-Sleep -Milliseconds 500
        $process.StandardInput.WriteLine($Password)
        
        $process.WaitForExit(30000) # 30秒超时
        
        if ($process.ExitCode -eq 0) {
            Write-Host "✅ 成功: $Description" -ForegroundColor Green
            return $true
        } else {
            Write-Host "❌ 失败: $Description" -ForegroundColor Red
            return $false
        }
    }
    catch {
        Write-Host "❌ 异常: $Description - $($_.Exception.Message)" -ForegroundColor Red
        return $false
    }
}

# 第1步：验证本地文件
Write-Host "🔍 第1步: 验证本地文件" -ForegroundColor Magenta
$requiredFiles = @(
    "bin\rigel-proxy-linux",
    "config\proxy1.yaml",
    "config\proxy2.yaml", 
    "config\target.yaml"
)

$allFilesExist = $true
foreach ($file in $requiredFiles) {
    if (Test-Path $file) {
        Write-Host "  ✅ $file" -ForegroundColor Green
    } else {
        Write-Host "  ❌ $file (缺失)" -ForegroundColor Red
        $allFilesExist = $false
    }
}

if (-not $allFilesExist) {
    Write-Host "❌ 缺少必要文件，请先编译程序和创建配置文件" -ForegroundColor Red
    exit 1
}

# 第2步：测试连接
Write-Host "`n🔗 第2步: 测试SSH连接" -ForegroundColor Magenta
foreach ($server in $servers.GetEnumerator()) {
    $result = Invoke-SSHCommandSafe -Server $server.Value.IP -Command "echo 'Connection test successful'" -Description "测试连接到$($server.Key)"
    if (-not $result) {
        Write-Host "❌ 无法连接到 $($server.Key)，停止部署" -ForegroundColor Red
        exit 1
    }
}

# 第3步：初始化环境
Write-Host "`n🔧 第3步: 初始化服务器环境" -ForegroundColor Magenta
$initCommands = @(
    "mkdir -p /opt/rigel/{bin,config,logs}",
    "apt update > /dev/null 2>&1",
    "apt install -y ufw curl wget > /dev/null 2>&1",
    "ufw --force enable > /dev/null 2>&1",
    "ufw allow ssh > /dev/null 2>&1",
    "ufw allow 8080/tcp > /dev/null 2>&1"
)

foreach ($server in $servers.GetEnumerator()) {
    Write-Host "初始化 $($server.Key)..." -ForegroundColor Yellow
    foreach ($cmd in $initCommands) {
        Invoke-SSHCommandSafe -Server $server.Value.IP -Command $cmd -Description "执行: $cmd" | Out-Null
    }
}

# 第4步：部署程序和配置
Write-Host "`n📦 第4步: 部署程序和配置文件" -ForegroundColor Magenta
foreach ($server in $servers.GetEnumerator()) {
    $serverName = $server.Key
    $serverIP = $server.Value.IP
    $configFile = $server.Value.Config
    
    Write-Host "部署到 $serverName ($serverIP)..." -ForegroundColor Yellow
    
    # 上传程序
    Copy-FileToServerSafe -Server $serverIP -LocalPath "bin\rigel-proxy-linux" -RemotePath "/opt/rigel/bin/rigel-proxy" -Description "上传程序到$serverName"
    
    # 上传配置
    Copy-FileToServerSafe -Server $serverIP -LocalPath "config\$configFile" -RemotePath "/opt/rigel/config.yaml" -Description "上传配置到$serverName"
    
    # 设置权限
    Invoke-SSHCommandSafe -Server $serverIP -Command "chmod +x /opt/rigel/bin/rigel-proxy" -Description "设置$serverName权限" | Out-Null
}

# 第5步：启动服务
Write-Host "`n🚀 第5步: 启动代理服务" -ForegroundColor Magenta
foreach ($server in $servers.GetEnumerator()) {
    $serverName = $server.Key
    $serverIP = $server.Value.IP
    
    $startCommand = "cd /opt/rigel && pkill rigel-proxy; nohup ./bin/rigel-proxy -c config.yaml > logs/rigel.log 2>&1 & echo 'Started'"
    Invoke-SSHCommandSafe -Server $serverIP -Command $startCommand -Description "启动$serverName服务"
}

# 第6步：等待启动
Write-Host "`n⏳ 等待服务启动..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# 第7步：验证部署
Write-Host "`n🔍 第7步: 验证部署状态" -ForegroundColor Magenta
foreach ($server in $servers.GetEnumerator()) {
    $serverName = $server.Key
    $serverIP = $server.Value.IP
    
    Write-Host "验证 $serverName ($serverIP):" -ForegroundColor Yellow
    
    # 测试端口连通性
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
}

Write-Host "`n🎉 部署完成！" -ForegroundColor Green
Write-Host "📋 下一步: 运行测试客户端" -ForegroundColor Yellow
Write-Host "  go run deploy\real-world-client.go" -ForegroundColor Gray
