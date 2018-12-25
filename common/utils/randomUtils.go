package utils

import (
	"math/rand"
	"time"
)

// GenerateRangeNum 生成固定范围的随机整数
func GenerateRangeNum(min, max int64) int64 {
	rand.Seed(time.Now().Unix())
	randNum := rand.Int63n(max-min) + min
	return randNum
}

// GetRandBarcode 根据日期生成随机条形码(int64)
func GetRandBarcode() int64 {
	return ToInt64(GetRandBarcodeS())
}

// GetRandBarcodeS 根据日期生成随机条形码(string)
func GetRandBarcodeS() string {
	num := GenerateRangeNum(10000, 99999)
	return GetDate() + Int64ToString(num)
}
