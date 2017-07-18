package comm

import (
	"fmt"
	"strconv"
)

func I32toa(n int32) string {
	return strconv.FormatInt(int64(n), 10)
}

func Atoi32(v string) int32 {
	n, err := strconv.ParseInt(v, 10, 32)
	if err == nil {
		return int32(n)
	} else {
		return 0
	}
}

func I64toa(n int64) string {
	return strconv.FormatInt(n, 10)
}

func Atoi64(v string) int64 {
	n, err := strconv.ParseInt(v, 10, 64)
	if err == nil {
		return n
	} else {
		return 0
	}
}

func Panic(v ...interface{}) {
	panic(fmt.Sprintln(v...))
}
