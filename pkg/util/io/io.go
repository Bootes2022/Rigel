package io

import (
	"io"
	"net"
	"sync"
)

// Join 在两个连接之间双向转发数据
func Join(c1, c2 net.Conn, bufferSize int) (inCount, outCount int64, err error) {
	var wait sync.WaitGroup
	var errOnce sync.Once
	var firstErr error

	wait.Add(2)

	// c1 -> c2
	go func() {
		defer wait.Done()
		var err error
		outCount, err = copyBuffer(c2, c1, bufferSize)
		if err != nil {
			errOnce.Do(func() {
				firstErr = err
			})
			c1.Close()
			c2.Close()
		}
	}()

	// c2 -> c1
	go func() {
		defer wait.Done()
		var err error
		inCount, err = copyBuffer(c1, c2, bufferSize)
		if err != nil {
			errOnce.Do(func() {
				firstErr = err
			})
			c1.Close()
			c2.Close()
		}
	}()

	wait.Wait()
	return inCount, outCount, firstErr
}

// copyBuffer 使用指定大小的缓冲区复制数据
func copyBuffer(dst io.Writer, src io.Reader, bufferSize int) (written int64, err error) {
	buf := make([]byte, bufferSize)
	return io.CopyBuffer(dst, src, buf)
}
