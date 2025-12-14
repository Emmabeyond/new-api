package common

import (
	"container/list"
	"sync"
	"sync/atomic"
	"time"
)

// CacheEntry 缓存条目
type CacheEntry[T any] struct {
	Key      string
	Value    T
	ExpireAt time.Time
	IsEmpty  bool // 空结果标记，防止缓存穿透
}

// IsExpired 检查条目是否过期
func (e *CacheEntry[T]) IsExpired() bool {
	if e.ExpireAt.IsZero() {
		return false
	}
	return time.Now().After(e.ExpireAt)
}

// LRUCache L1 本地 LRU 缓存
// 线程安全，支持泛型
type LRUCache[T any] struct {
	capacity  int
	items     map[string]*list.Element
	evictList *list.List
	mu        sync.RWMutex

	// 统计指标
	hits   int64
	misses int64
}

// NewLRUCache 创建新的 LRU 缓存
func NewLRUCache[T any](capacity int) *LRUCache[T] {
	if capacity <= 0 {
		capacity = 10000 // 默认容量
	}
	return &LRUCache[T]{
		capacity:  capacity,
		items:     make(map[string]*list.Element),
		evictList: list.New(),
	}
}

// Get 获取缓存值
// 返回值、是否命中、是否为空结果标记
func (c *LRUCache[T]) Get(key string) (T, bool, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zero T

	elem, ok := c.items[key]
	if !ok {
		atomic.AddInt64(&c.misses, 1)
		return zero, false, false
	}

	entry := elem.Value.(*CacheEntry[T])

	// 检查是否过期
	if entry.IsExpired() {
		c.removeElement(elem)
		atomic.AddInt64(&c.misses, 1)
		return zero, false, false
	}

	// 移动到链表头部（最近使用）
	c.evictList.MoveToFront(elem)
	atomic.AddInt64(&c.hits, 1)

	return entry.Value, true, entry.IsEmpty
}

// Set 设置缓存值
func (c *LRUCache[T]) Set(key string, value T, ttl time.Duration) {
	c.SetWithEmpty(key, value, ttl, false)
}

// SetEmpty 设置空结果缓存（防止缓存穿透）
func (c *LRUCache[T]) SetEmpty(key string, ttl time.Duration) {
	var zero T
	c.SetWithEmpty(key, zero, ttl, true)
}

// SetWithEmpty 设置缓存值，支持空结果标记
func (c *LRUCache[T]) SetWithEmpty(key string, value T, ttl time.Duration, isEmpty bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 计算过期时间
	var expireAt time.Time
	if ttl > 0 {
		expireAt = time.Now().Add(ttl)
	}

	// 如果 key 已存在，更新值并移动到头部
	if elem, ok := c.items[key]; ok {
		c.evictList.MoveToFront(elem)
		entry := elem.Value.(*CacheEntry[T])
		entry.Value = value
		entry.ExpireAt = expireAt
		entry.IsEmpty = isEmpty
		return
	}

	// 检查容量，需要时淘汰最久未使用的条目
	for c.evictList.Len() >= c.capacity {
		c.evictOldest()
	}

	// 添加新条目到头部
	entry := &CacheEntry[T]{
		Key:      key,
		Value:    value,
		ExpireAt: expireAt,
		IsEmpty:  isEmpty,
	}
	elem := c.evictList.PushFront(entry)
	c.items[key] = elem
}

// Delete 删除缓存条目
func (c *LRUCache[T]) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.removeElement(elem)
	}
}

// Len 返回当前缓存条目数
func (c *LRUCache[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.evictList.Len()
}

// Capacity 返回缓存容量
func (c *LRUCache[T]) Capacity() int {
	return c.capacity
}

// Stats 返回缓存统计信息
func (c *LRUCache[T]) Stats() (hits, misses int64, size, capacity int) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return atomic.LoadInt64(&c.hits), atomic.LoadInt64(&c.misses), c.evictList.Len(), c.capacity
}

// HitRate 返回缓存命中率
func (c *LRUCache[T]) HitRate() float64 {
	hits := atomic.LoadInt64(&c.hits)
	misses := atomic.LoadInt64(&c.misses)
	total := hits + misses
	if total == 0 {
		return 0
	}
	return float64(hits) / float64(total)
}

// Clear 清空缓存
func (c *LRUCache[T]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element)
	c.evictList.Init()
}

// CleanExpired 清理过期条目
// 返回清理的条目数
func (c *LRUCache[T]) CleanExpired() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	count := 0
	now := time.Now()

	// 从尾部开始检查（最久未使用的条目更可能过期）
	for elem := c.evictList.Back(); elem != nil; {
		entry := elem.Value.(*CacheEntry[T])
		prev := elem.Prev()

		if !entry.ExpireAt.IsZero() && now.After(entry.ExpireAt) {
			c.removeElement(elem)
			count++
		}

		elem = prev
	}

	return count
}

// evictOldest 淘汰最久未使用的条目
func (c *LRUCache[T]) evictOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.removeElement(elem)
	}
}

// removeElement 移除元素
func (c *LRUCache[T]) removeElement(elem *list.Element) {
	c.evictList.Remove(elem)
	entry := elem.Value.(*CacheEntry[T])
	delete(c.items, entry.Key)
}
