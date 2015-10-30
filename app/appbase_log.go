package app

import (
	"fmt"
	"strings"
)

const (
	LogDebugLevel = 0                // 日志等级 : 调试信息
	LogInfoLevel  = 1                // 日志等级 : 普通信息
	LogWarnLevel  = 2                // 日志等级 : 警告信息
	LogErrorLevel = 3                // 日志等级 : 错误信息
	LogFatalLevel = 4                // 日志等级 : 致命信息
	LogMaxLevel   = 5                // 日志最大等级
	LogLimitLevel = LogInfoLevel     // 显示这个等级之上的日志(控制台)
	LogBuffMax    = 20 * 1024 * 1024 // 日志缓冲
)

const (
	AppName     = "##[AppName]##"
	LogBuffSize = 10 * 1024 * 1024
	LogDir      = "../log"
	ProfFile    = AppName + "_prof.log"
	LogFileName = LogDir + "/" + AppName + ".log"
)

// 线程日志 : 生成日志头
func (this *AppBase) MakeLogHeader() {
	this.log_Header[LogDebugLevel] = this.log_TimeString + " [D] "
	this.log_Header[LogInfoLevel] = this.log_TimeString + " [I] "
	this.log_Header[LogWarnLevel] = this.log_TimeString + " [W] "
	this.log_Header[LogErrorLevel] = this.log_TimeString + " [E] "
	this.log_Header[LogFatalLevel] = this.log_TimeString + " [F] "
}

// 线程日志 : 调试[D]级别日志
func (this *AppBase) LogDebug(f string, v ...interface{}) {
	this.LogBase(LogDebugLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 信息[I]级别日志
func (this *AppBase) LogInfo(f string, v ...interface{}) {
	this.LogBase(LogInfoLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 警告[W]级别日志
func (this *AppBase) LogWarn(f string, v ...interface{}) {
	this.LogBase(LogWarnLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 错误[E]级别日志
func (this *AppBase) LogError(f string, v ...interface{}) {
	this.LogBase(LogErrorLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 致命[F]级别日志
func (this *AppBase) LogFatal(f string, v ...interface{}) {
	this.LogBase(LogFatalLevel, fmt.Sprintf(f, v...))
}

// 线程日志 : 手动分级日志
func (this *AppBase) LogBase(level int, info string) {
	if level >= LogDebugLevel && level < LogMaxLevel {
		s := this.log_Header[level] + info
		s = strings.Replace(s, "\n", "\n"+this.log_Header[level], -1) + "\n"

		this.log_FileBuff.WriteString(s)

		if level >= LogLimitLevel {
			fmt.Print(s)
		}
	} else {
		fmt.Println("LogBase : level failed : ", level)
	}
}
