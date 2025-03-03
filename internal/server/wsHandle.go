package server

import (
	//"github.com/go-redis/redis"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/websocket"
)

type ClientManager struct {
	clients     map[*Client]bool
	broadcast   chan []byte
	register    chan *Client
	unregister  chan *Client
	redisPubSub *redis.PubSub
}

type Client struct {
	uuid   string
	conn   *websocket.Conn
	send   chan []byte
	topics map[string]bool // 订阅的频道
}
