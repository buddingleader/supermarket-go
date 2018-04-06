package common

import (
	"strconv"
	"supermarket-go/log"
)

// ToInt64 string转换为int64，转换失败则为0
func ToInt64(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		log.ErrorLog(err)
		return 0
	}
	return v
}

// Int64ToString int64转换为string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Float64ToString Float64转换为string，保留两位小数
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}
