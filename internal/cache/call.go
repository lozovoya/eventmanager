package cache

import (
	"EventManager/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"sync"
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
	wg := sync.WaitGroup{}
	resultCh := make(chan []*model.Call)
	keysCh := make(chan []string)

	go func() {
		for r := range resultCh {
			log.Printf("recording calls to slice: %d", len(calls))
			result := r
			calls = append(calls, result...)
		}
		log.Println("recording calls is finished")

	}()

	go func() {
		var cursor = 0
		var counter = 1000
		for {
			log.Printf("redis works, counter %d", counter)
			data, err := redis.Values(redis.DoWithTimeout(conn, time.Millisecond*2000, "SCAN", cursor, "COUNT", counter))
			if err != nil {
				fmt.Errorf("GetCallsSnapshot %w", err)
				close(keysCh)
			}
			cursor, _ = redis.Int(data[0], nil)
			keys, _ := redis.Strings(data[1], nil)
			log.Printf("keys length: %d", len(keys))
			keysCh <- keys
			if cursor == 0 {
				close(keysCh)
				break
			}
		}
		log.Println("redis done")
	}()

	i := 0
	for keys := range keysCh {
		wg.Add(1)
		i++
		log.Printf("start gr %d", i)
		go func(keys []string, wg *sync.WaitGroup) {
			conn, err := c.pool.GetContext(ctx)
			if err != nil {
				fmt.Errorf("GetCallsSnapshot %w", err)
			}
			defer conn.Close()
			var calls []*model.Call
			for _, key := range keys {
				//log.Printf("searching keys in redis: %s", key)
				var call model.Call
				data, _ := redis.Bytes(redis.DoWithTimeout(conn, time.Millisecond*2000, "GET", key))
				if err != nil {
					fmt.Errorf("GetCallsSnapshot %w", err)
				}
				err = json.Unmarshal(data, &call)
				if err != nil {
					fmt.Errorf("GetCallsSnapshot %w", err)
				}
				calls = append(calls, &call)
			}
			log.Println("send result to channel")
			resultCh <- calls
			wg.Done()
		}(keys, &wg)
	}

	log.Println("wait for gr")
	wg.Wait()
	close(resultCh)

	return calls, nil
}
