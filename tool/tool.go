package tool

import (
	"flag"
	"unsafe"
)

func IsDebug() bool {
	return flag.Lookup("test.v") != nil
}

const (
	FormatIso8601Date                   = "2006-01-02"
	FORMAT_ISO8601_DATE_TIME            = "2006-01-02 15:04:05"
	FORMAT_ISO8601_DATE_TIME_MILLI      = "2006-01-02 15:04:05.000"
	FORMAT_ISO8601_DATE_TIME_MILLI_ZONE = "2006-01-02 15:04:05.000Z07:00"
	FORMAT_ISO8601_DATE_TIME_MICRO      = "2006-01-02 15:04:05.000000"
	FORMAT_ISO8601_DATE_TIME_MICRO_ZONE = "2006-01-02 15:04:05.000000Z07:00"
	FORMAT_ISO8601_DATE_TIME_NANO       = "2006-01-02 15:04:05.000000000"
	FORMAT_ISO8601_DATE_TIME_NANO_ZONE  = "2006-01-02 15:04:05.00000000007:00"
)

// type = struct string {
//
//    uint8 *str;
//
//    int len;
//
//}
//
//(gdb) ptype b
//
//type = struct []uint8 {
//
//    uint8 *array;
//
//    int len;
//
//    int cap;
//
//}
// go build -gcflags=-m 未发生逃逸
func StrToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}
