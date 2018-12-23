package main

import (
	"github.com/wangff15386/supermarket-go/log"
)

func main() {
	myLogger := log.GetLogger("tet", "")
	myLogger.Println("This is an info message, with colors (if the output is terminal)")
}
