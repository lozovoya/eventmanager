package cache

import (
	"context"
	"github.com/gomodule/redigo/redis"
)

const TTL = 120

func InitCache(addr string) *redis.Pool {
	myCache := &redis.Pool{
		DialContext: func(ctx context.Context) (redis.Conn, error) {
			return redis.DialURL(addr)
		},
	}
	return myCache
}




