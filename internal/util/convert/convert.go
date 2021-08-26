package convert

import (
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"
)

func AtoI64(src string) int64 {
	num, _ := strconv.ParseInt(src, 10, 64)
	return num
}

func I64toA(src int64) string {
	return strconv.FormatInt(src, 10)
}


func StringToByte(src string) []byte {
	str := (*reflect.StringHeader)(unsafe.Pointer(&src))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: str.Data, Len: str.Len, Cap: str.Len}))
}

func ByteToString(src []byte) string {
	str := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: str.Data, Len: str.Len}))
}

func MustMarshal(v interface{}) string {
	bytes, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return ByteToString(bytes)
}