package convert

import "strconv"

func AtoI64(src string) int64 {
	num, _ := strconv.ParseInt(src, 10, 64)
	return num
}

func I64toA(src int64) string {
	return strconv.FormatInt(src, 10)
}