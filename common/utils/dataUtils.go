package utils

import (
	"bytes"
	"time"
)

// GetDate  返回20060102格式数据
func GetDate() string {
	return time.Now().Format("20060102")
}

// GetDate1  返回2006-01-02格式数据
func GetDate1() string {
	return time.Now().Format("2006-01-02")
}

// GetTime  返回150405格式数据
func GetTime() string {
	return time.Now().Format("150405")
}

// GetTime1  返回15:04:05格式数据
func GetTime1() string {
	return time.Now().Format("15:04:05")
}

// GetTimeStamp 返回20060102150405格式数据
func GetTimeStamp() string {
	sb := bytes.Buffer{}
	sb.WriteString(GetDate())
	sb.WriteString(GetTime())
	return sb.String()
}

// GetTimeStamp1 返回2006-01-02 15:04:05格式数据
func GetTimeStamp1() string {
	sb := bytes.Buffer{}
	sb.WriteString(GetDate1())
	sb.WriteString(" ")
	sb.WriteString(GetTime1())
	return sb.String()
}
