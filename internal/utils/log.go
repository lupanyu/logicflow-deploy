package utils

import (
	"fmt"
	"log"
	"path/filepath"
	"runtime"
)

func GetCallerInfo() string {
	_, file, line, _ := runtime.Caller(1) // +1 跳过自身调用
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func InitLog() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
}
