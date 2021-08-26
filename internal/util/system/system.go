package system

import (
	"log"
	"os"
)

/*
工具包
系统函数
*/

func Hostname() string {
	if hostname, err := os.Hostname(); err != nil {
		log.Panic(err)
		return ""
	} else {
		return hostname
	}
}
