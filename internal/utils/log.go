package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func GetCallerInfo(skip int) string {
	_, file, line, _ := runtime.Caller(1) // +1 跳过自身调用
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}
