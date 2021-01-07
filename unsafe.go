package utils

import "unsafe"

func UnsafeInterfaceByteSize(res interface{}) float32 {
	return float32(unsafe.Sizeof(res)+unsafe.Sizeof([1000]int32{})) / 1024
}
