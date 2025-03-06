package protocol

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"log"
	"logicflow-deploy/internal/utils"
)

var jwtKey = []byte("your_secret_key")

// 定义 Claims 结构体
type Claims struct {
	AgentID string `json:"agent_id"`
	jwt.StandardClaims
}

type MessageAuthResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 定义 doAuthentication 函数
func doAuthentication(conn *websocket.Conn) string {
	// 读取客户端发送的 JWT
	var tokenStr string
	if err := conn.ReadJSON(&tokenStr); err != nil {
		log.Printf(" [%s]读取 JWT 失败: %v", utils.GetCallerInfo(), err)
		return ""
	}

	// 解析 JWT
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		log.Printf(" [%s]解析 JWT 失败: %v", utils.GetCallerInfo(), err)
		return ""
	}

	if !tkn.Valid {
		log.Println("无效的 JWT")
		return ""
	}

	// 返回 Agent ID
	return claims.AgentID
}
