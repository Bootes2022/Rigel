# Target服务器修复脚本

param(
    [string]$Password = "20010221cc."
)

Write-Host "🔧 Target服务器修复脚本" -ForegroundColor Green
Write-Host "========================" -ForegroundColor Green

$targetIP = "120.27.113.148"

Write-Host "🎯 修复Target服务器: $targetIP" -ForegroundColor Yellow
Write-Host ""

# 函数：执行SSH命令
function Invoke-TargetCommand {
    param(
        [string]$Command,
        [string]$Description
    )
    
    Write-Host "🔧 $Description..." -ForegroundColor Yellow
    
    try {
        $result = ssh -o StrictHostKeyChecking=no -o ConnectTimeout=10 root@$targetIP $Command 2>&1
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ 成功: $Description" -ForegroundColor Green
            return $result
        } else {
            Write-Host "❌ 失败: $Description" -ForegroundColor Red
            Write-Host "错误: $result" -ForegroundColor Red
            return $null
        }
    }
    catch {
        Write-Host "❌ 异常: $Description - $($_.Exception.Message)" -ForegroundColor Red
        return $null
    }
}

# 第1步：检查当前状态
Write-Host "🔍 第1步: 检查当前状态" -ForegroundColor Magenta

# 检查网络连通性
$ping = Test-Connection -ComputerName $targetIP -Count 1 -Quiet
if ($ping) {
    Write-Host "✅ 网络连通正常" -ForegroundColor Green
} else {
    Write-Host "❌ 网络连通失败，无法继续" -ForegroundColor Red
    exit 1
}

# 检查SSH连接
$sshTest = Invoke-TargetCommand "echo 'SSH connection successful'" "测试SSH连接"
if (-not $sshTest) {
    Write-Host "❌ SSH连接失败，无法继续" -ForegroundColor Red
    exit 1
}

# 第2步：停止现有进程
Write-Host "`n🛑 第2步: 停止现有进程" -ForegroundColor Magenta
Invoke-TargetCommand "pkill rigel-proxy" "停止rigel-proxy进程"

# 第3步：检查文件状态
Write-Host "`n📁 第3步: 检查文件状态" -ForegroundColor Magenta
$fileCheck = Invoke-TargetCommand "ls -la /opt/rigel/" "检查目录结构"
if ($fileCheck) {
    Write-Host $fileCheck -ForegroundColor Gray
}

$binaryCheck = Invoke-TargetCommand "ls -la /opt/rigel/bin/rigel-proxy" "检查程序文件"
if ($binaryCheck) {
    Write-Host $binaryCheck -ForegroundColor Gray
}

$configCheck = Invoke-TargetCommand "ls -la /opt/rigel/config.yaml" "检查配置文件"
if ($configCheck) {
    Write-Host $configCheck -ForegroundColor Gray
}

# 第4步：重新上传文件（如果需要）
Write-Host "`n📤 第4步: 重新上传文件" -ForegroundColor Magenta

# 上传程序
Write-Host "上传程序文件..." -ForegroundColor Yellow
$uploadResult = scp -o StrictHostKeyChecking=no config\target.yaml root@${targetIP}:/opt/rigel/config.yaml 2>&1
if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ 配置文件上传成功" -ForegroundColor Green
} else {
    Write-Host "❌ 配置文件上传失败: $uploadResult" -ForegroundColor Red
}

# 第5步：设置权限
Write-Host "`n🔐 第5步: 设置权限" -ForegroundColor Magenta
Invoke-TargetCommand "chmod +x /opt/rigel/bin/rigel-proxy" "设置程序权限"

# 第6步：检查配置文件
Write-Host "`n📋 第6步: 检查配置文件" -ForegroundColor Magenta
$configContent = Invoke-TargetCommand "cat /opt/rigel/config.yaml" "读取配置文件"
if ($configContent) {
    Write-Host "配置文件内容:" -ForegroundColor Gray
    Write-Host $configContent -ForegroundColor Gray
}

# 第7步：启动服务
Write-Host "`n🚀 第7步: 启动服务" -ForegroundColor Magenta
$startCommand = "cd /opt/rigel && nohup ./bin/rigel-proxy -c config.yaml > logs/rigel.log 2>&1 & echo 'Target started'"
Invoke-TargetCommand $startCommand "启动Target服务"

# 第8步：等待启动
Write-Host "`n⏳ 等待服务启动..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

# 第9步：验证启动
Write-Host "`n🔍 第9步: 验证启动状态" -ForegroundColor Magenta

# 检查进程
$processCheck = Invoke-TargetCommand "ps aux | grep rigel-proxy | grep -v grep" "检查进程状态"
if ($processCheck) {
    Write-Host "进程状态:" -ForegroundColor Gray
    Write-Host $processCheck -ForegroundColor Gray
}

# 检查端口
$portCheck = Invoke-TargetCommand "netstat -tlnp | grep :8080" "检查端口监听"
if ($portCheck) {
    Write-Host "端口状态:" -ForegroundColor Gray
    Write-Host $portCheck -ForegroundColor Gray
}

# 检查日志
$logCheck = Invoke-TargetCommand "tail -5 /opt/rigel/logs/rigel.log" "检查日志"
if ($logCheck) {
    Write-Host "最新日志:" -ForegroundColor Gray
    Write-Host $logCheck -ForegroundColor Gray
}

# 第10步：测试连通性
Write-Host "`n🧪 第10步: 测试连通性" -ForegroundColor Magenta
try {
    $connection = Test-NetConnection -ComputerName $targetIP -Port 8080 -WarningAction SilentlyContinue
    if ($connection.TcpTestSucceeded) {
        Write-Host "✅ Target端口8080连通正常" -ForegroundColor Green
    } else {
        Write-Host "❌ Target端口8080连接失败" -ForegroundColor Red
        Write-Host "可能需要检查防火墙设置或程序启动状态" -ForegroundColor Yellow
    }
}
catch {
    Write-Host "❌ 连通性测试异常: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "🎯 Target修复完成！" -ForegroundColor Green
