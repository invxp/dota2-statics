package system

import (
	"fmt"
	"os"
)

/*
工具包
系统函数
*/

func Hostname() string {
	if hostname, err := os.Hostname(); err != nil {
		return fmt.Sprintf("unknown-%v",err)
	} else {
		return hostname
	}
}
