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
