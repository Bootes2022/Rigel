package net

import (
	"fmt"
	"net"
	"time"
)

// IsPortAvailable 检查端口是否可用
func IsPortAvailable(port int) bool {
	// 占位符实现，实际使用时需要完善
	addr := fmt.Sprintf(":%d", port)
	conn, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// FindAvailablePort 查找可用端口
func FindAvailablePort(startPort int) (int, error) {
	// 占位符实现，实际使用时需要完善
	for port := startPort; port < startPort+1000; port++ {
		if IsPortAvailable(port) {
			return port, nil
		}
	}
	return 0, fmt.Errorf("no available port found in range %d-%d", startPort, startPort+1000)
}

// SetKeepAlive 设置 TCP 连接的 keepalive
func SetKeepAlive(conn net.Conn, keepalive bool, keepaliveIdle time.Duration) error {
	// 占位符实现，实际使用时需要完善
	tcpConn, ok := conn.(*net.TCPConn)
	if !ok {
		return fmt.Errorf("not a TCP connection")
	}

	if err := tcpConn.SetKeepAlive(keepalive); err != nil {
		return err
	}

	if keepalive {
		if err := tcpConn.SetKeepAlivePeriod(keepaliveIdle); err != nil {
			return err
		}
	}

	return nil
}

// ParseAddress 解析地址字符串
func ParseAddress(addr string, defaultPort int) (string, int, error) {
	// 占位符实现，实际使用时需要完善
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		// 可能没有指定端口
		return addr, defaultPort, nil
	}

	portNum := defaultPort
	if port != "" {
		var err error
		fmt.Sscanf(port, "%d", &portNum)
		if err != nil {
			return "", 0, fmt.Errorf("invalid port: %s", port)
		}
	}

	return host, portNum, nil
}
