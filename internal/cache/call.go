package cache

import (
	"EventManager/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"time"
)

type callCache struct {
	pool *redis.Pool
}

func NewCallCache(pool *redis.Pool) Call {
	return &callCache{pool: pool}
}

func (c callCache) CallToCache(ctx context.Context, call *model.Call) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("CallToCache %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(err)
		}
	}()
	key := fmt.Sprintf("%s:%s", call.Queue_ID, call.CallID)
	value, err := json.Marshal(call)
	_, err = redis.DoWithTimeout(conn, time.Millisecond*2000, "SET", key, value)
	if err != nil {
		return fmt.Errorf("CallToCache %w", err)
	}
	return nil
}

func (c callCache) CallFromCache(ctx context.Context, queueID, callID string) error {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("CallFromCache %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(err)
		}
	}()
	key := fmt.Sprintf("%s:%s", queueID, callID)
	_, err = redis.DoWithTimeout(conn, time.Millisecond*2000, "DEL", key)
	if err != nil {
		return fmt.Errorf("CallToCache %w", err)
	}
	return nil
}

func (c callCache) GetCallsSnapshot(ctx context.Context) ([]*model.Call, error) {
	conn, err := c.pool.GetContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("GetCallsSnapshot %w", err)
	}
	defer func() {
		if cerr := conn.Close(); cerr != nil {
			log.Println(err)
		}
	}()
	var calls []*model.Call
	var cursor = 0
	for {
		data, err := redis.Values(redis.DoWithTimeout(conn, time.Millisecond*2000, "SCAN", cursor))
		if err != nil {
			return nil, fmt.Errorf("GetCallsSnapshot %w", err)
		}
		cursor, _ = redis.Int(data[0], nil)
		keys, _ := redis.Strings(data[1], nil)

		for _, key := range keys {
			var call model.Call
			data, err := redis.Bytes(redis.DoWithTimeout(conn, time.Millisecond*2000, "GET", key))
			if err != nil {
				return nil, fmt.Errorf("GetCallsSnapshot %w", err)
			}
			err = json.Unmarshal(data, &call)
			if err != nil {
				return nil, fmt.Errorf("GetCallsSnapshot %w", err)
			}
			calls = append(calls, &call)
		}
		if cursor == 0 {
			break
		}
		log.Println("redis done")
	}
	return calls, nil
}
