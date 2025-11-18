package logger

import (
	"log"
	"os"
)

// Logger 封装标准日志
type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	error *log.Logger
	debug *log.Logger
}

// New 创建新的Logger实例
func New() *Logger {
	return &Logger{
		info:  log.New(os.Stdout, "[INFO] ", log.LstdFlags|log.Lshortfile),
		warn:  log.New(os.Stdout, "[WARN] ", log.LstdFlags|log.Lshortfile),
		error: log.New(os.Stderr, "[ERROR] ", log.LstdFlags|log.Lshortfile),
		debug: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags|log.Lshortfile),
	}
}

// Info 输出信息日志
func (l *Logger) Info(v ...interface{}) {
	l.info.Println(v...)
}

// Warn 输出警告日志
func (l *Logger) Warn(v ...interface{}) {
	l.warn.Println(v...)
}

// Error 输出错误日志
func (l *Logger) Error(v ...interface{}) {
	l.error.Println(v...)
}

// Debug 输出调试日志
func (l *Logger) Debug(v ...interface{}) {
	l.debug.Println(v...)
}

// Fatal 输出致命错误并退出
func (l *Logger) Fatal(v ...interface{}) {
	l.error.Fatal(v...)
}

// 全局logger实例
var defaultLogger = New()

// Info 全局Info方法
func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

// Warn 全局Warn方法
func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

// Error 全局Error方法
func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

// Debug 全局Debug方法
func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

// Fatal 全局Fatal方法
func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}
