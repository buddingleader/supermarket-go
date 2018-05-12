package log

import (
	"fmt"
	"log"
	"os"
)

// 日志存放路径
const (
	LOGPATH = "/supermarket/logs/"
)

// 日志记录
var (
	infoLog  *log.Logger
	errorLog *log.Logger
)

func init() {
	var err error
	if infoLog, err = initLogFile("xxx_Info.log", "[Info]"); err != nil {
		log.Panic("创建Info日志失败!")
	}
	if errorLog, err = initLogFile("xxx_Error.log", "[Error]"); err != nil {
		log.Panic("创建Error日志失败!")
	}
}

func initLogFile(fileName string, level string) (*log.Logger, error) {
	logFile, err := os.OpenFile(fileName, os.O_APPEND, 0666) //打开文件
	if err != nil {
		logFile, err = os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666) //打开文件
	}
	defer logFile.Close()
	newLog := log.New(logFile, level, log.LstdFlags|log.Llongfile)
	newLog.Printf("A %s message here", level)
	return newLog, nil
}

// ErrorLog 记录错误日志
func ErrorLog(message ...interface{}) {
	fmt.Println(message)
	errorLog.Println(message)
}

// InfoLog 记录普通日志
func InfoLog(message ...interface{}) {
	fmt.Println(message)
	infoLog.Println(message)
}
