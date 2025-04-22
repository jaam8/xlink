package pool

import (
	"google.golang.org/grpc"
	"time"
)

type Config struct {
	Address        string `json:"address"`
	MaxConnections int    `json:"max_connections"`
	MinConnections int    `json:"min_connections"`
	DialOptions    []grpc.DialOption
	DialTimeout    time.Duration
}
