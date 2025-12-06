package common

import (
	"hash/fnv"
	"sync"
	"time"
)

// 分片数量，用于减少锁竞争
const shardCount = 32

type rateLimiterShard struct {
	store map[string]*[]int64
	mutex sync.Mutex
}

type InMemoryRateLimiter struct {
	shards             [shardCount]*rateLimiterShard
	expirationDuration time.Duration
	initialized        bool
	initMutex          sync.Mutex
}

// getShard 根据key获取对应的分片
func (l *InMemoryRateLimiter) getShard(key string) *rateLimiterShard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return l.shards[h.Sum32()%shardCount]
}

func (l *InMemoryRateLimiter) Init(expirationDuration time.Duration) {
	if l.initialized {
		return
	}
	l.initMutex.Lock()
	defer l.initMutex.Unlock()
	if l.initialized {
		return
	}
	for i := 0; i < shardCount; i++ {
		l.shards[i] = &rateLimiterShard{
			store: make(map[string]*[]int64),
		}
	}
	l.expirationDuration = expirationDuration
	l.initialized = true
	if expirationDuration > 0 {
		go l.clearExpiredItems()
	}
}

func (l *InMemoryRateLimiter) clearExpiredItems() {
	for {
		time.Sleep(l.expirationDuration)
		now := time.Now().Unix()
		for i := 0; i < shardCount; i++ {
			shard := l.shards[i]
			shard.mutex.Lock()
			for key := range shard.store {
				queue := shard.store[key]
				size := len(*queue)
				if size == 0 || now-(*queue)[size-1] > int64(l.expirationDuration.Seconds()) {
					delete(shard.store, key)
				}
			}
			shard.mutex.Unlock()
		}
	}
}

// Request parameter duration's unit is seconds
func (l *InMemoryRateLimiter) Request(key string, maxRequestNum int, duration int64) bool {
	shard := l.getShard(key)
	shard.mutex.Lock()
	defer shard.mutex.Unlock()
	// [old <-- new]
	queue, ok := shard.store[key]
	now := time.Now().Unix()
	if ok {
		if len(*queue) < maxRequestNum {
			*queue = append(*queue, now)
			return true
		} else {
			if now-(*queue)[0] >= duration {
				*queue = (*queue)[1:]
				*queue = append(*queue, now)
				return true
			} else {
				return false
			}
		}
	} else {
		s := make([]int64, 0, maxRequestNum)
		shard.store[key] = &s
		*(shard.store[key]) = append(*(shard.store[key]), now)
	}
	return true
}
