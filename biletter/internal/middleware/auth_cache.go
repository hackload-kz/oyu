package middleware

import (
	"biletter/internal/domain/entity"
	"sync"
	"time"
)

type authCacheEntry struct {
	user entity.AuthUser
	exp  int64
	neg  bool // негативный кэш
}

type authCache struct {
	mu     sync.RWMutex
	data   map[string]authCacheEntry
	ttlOK  time.Duration
	ttlBad time.Duration
	max    int // простой "poor man's" ограничитель, сбрасываем при переполнении
}

func newAuthCache() *authCache {
	return &authCache{
		data:   make(map[string]authCacheEntry, 8192),
		ttlOK:  120 * time.Second,
		ttlBad: 5 * time.Second,
		max:    20000,
	}
}

func (c *authCache) get(k string) (e authCacheEntry, ok bool) {
	now := time.Now().UnixNano()
	c.mu.RLock()
	e, ok = c.data[k]
	c.mu.RUnlock()
	if !ok || e.exp < now {
		return authCacheEntry{}, false
	}
	return e, true
}
func (c *authCache) setOK(key string, u entity.AuthUser) {
	c.mu.Lock()
	if len(c.data) > c.max {
		c.data = make(map[string]authCacheEntry, c.max) // грубая очистка
	}
	c.data[key] = authCacheEntry{user: u, exp: time.Now().Add(c.ttlOK).UnixNano()}
	c.mu.Unlock()
}
func (c *authCache) setBad(key string) {
	c.mu.Lock()
	if len(c.data) > c.max {
		c.data = make(map[string]authCacheEntry, c.max)
	}
	c.data[key] = authCacheEntry{neg: true, exp: time.Now().Add(c.ttlBad).UnixNano()}
	c.mu.Unlock()
}
