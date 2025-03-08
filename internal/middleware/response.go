package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"log"
)

// 自定义 ResponseWriter 用于捕获响应数据
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// 重写 Write 方法，将数据同时写入缓冲区和原始 ResponseWriter
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// 重写 WriteString 方法（处理直接写入字符串的情况）
func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 中间件：打印响应内容
func LogResponseBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 替换为自定义的 ResponseWriter
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// 处理请求
		c.Next()

		// 打印响应状态码和内容
		statusCode := c.Writer.Status()
		responseBody := blw.body.String()
		log.Printf("Response Status: %d | Body: %s\n", statusCode, responseBody)
	}
}
