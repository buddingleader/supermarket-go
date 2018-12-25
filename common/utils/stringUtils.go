package utils

import (
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

// ToInt64 convert string to int64
func ToInt64(str string) int64 {
	v, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		panic(err)
	}
	return v
}

// Int64ToString convert int64 to string
func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

// Float64ToString convert float64 to string, contains two
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

// ToString convert interface to string
func ToString(obj interface{}) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		return obj.(string)
	}

	return string(b)
}
