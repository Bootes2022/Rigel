package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

var (
	// 日志级别
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"

	// 日志实例
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger

	// 当前日志级别
	currentLevel string

	// 日志输出
	logWriter io.Writer
)

// Init 初始化日志
func Init(level, logFile string) error {
	// 占位符实现，实际使用时需要完善
	currentLevel = strings.ToLower(level)

	// 设置日志输出
	logWriter = os.Stdout
	if logFile != "" {
		file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file: %v", err)
		}
		logWriter = io.MultiWriter(os.Stdout, file)
	}

	// 创建日志实例
	debugLogger = log.New(logWriter, "[DEBUG] ", log.LstdFlags|log.Lshortfile)
	infoLogger = log.New(logWriter, "[INFO] ", log.LstdFlags)
	warnLogger = log.New(logWriter, "[WARN] ", log.LstdFlags)
	errorLogger = log.New(logWriter, "[ERROR] ", log.LstdFlags|log.Lshortfile)
	fatalLogger = log.New(logWriter, "[FATAL] ", log.LstdFlags|log.Lshortfile)

	return nil
}

// Debug 输出调试日志
func Debug(v ...interface{}) {
	if currentLevel == LevelDebug {
		debugLogger.Println(v...)
	}
}

// Debugf 输出格式化调试日志
func Debugf(format string, v ...interface{}) {
	if currentLevel == LevelDebug {
		debugLogger.Printf(format, v...)
	}
}

// Info 输出信息日志
func Info(v ...interface{}) {
	if currentLevel == LevelDebug || currentLevel == LevelInfo {
		infoLogger.Println(v...)
	}
}

// Infof 输出格式化信息日志
func Infof(format string, v ...interface{}) {
	if currentLevel == LevelDebug || currentLevel == LevelInfo {
		infoLogger.Printf(format, v...)
	}
}

// Warn 输出警告日志
func Warn(v ...interface{}) {
	if currentLevel != LevelError && currentLevel != LevelFatal {
		warnLogger.Println(v...)
	}
}

// Warnf 输出格式化警告日志
func Warnf(format string, v ...interface{}) {
	if currentLevel != LevelError && currentLevel != LevelFatal {
		warnLogger.Printf(format, v...)
	}
}

// Error 输出错误日志
func Error(v ...interface{}) {
	if currentLevel != LevelFatal {
		errorLogger.Println(v...)
	}
}

// Errorf 输出格式化错误日志
func Errorf(format string, v ...interface{}) {
	if currentLevel != LevelFatal {
		errorLogger.Printf(format, v...)
	}
}

// Fatal 输出致命错误日志并退出
func Fatal(v ...interface{}) {
	fatalLogger.Println(v...)
	os.Exit(1)
}

// Fatalf 输出格式化致命错误日志并退出
func Fatalf(format string, v ...interface{}) {
	fatalLogger.Printf(format, v...)
	os.Exit(1)
}
