package main

import (
	"github.com/wangff15386/supermarket-go/bussiness"
	"github.com/wangff15386/supermarket-go/conf"
	"github.com/wangff15386/supermarket-go/local/leveldb"
)

func main() {
	conf.Initial("../conf/app.conf")
	leveldb.Initial()
	bussiness.OpenBussiness()
}
