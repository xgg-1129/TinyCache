package TinyCache

import "sync"

type Cache struct {
	mu sync.Mutex
	cache *LRUCache
	maxBytes int64
}

func NewCache(cap int64)*Cache{
	c:=new(Cache)
	c.cache=NewLRUCache(cap)
	c.maxBytes=cap
	return c
}

func (c *Cache) get(key string)(Data,bool){
	c.mu.Lock()
	defer c.mu.Unlock()
	if value , ok := c.cache.Get(key);ok{
		return value.(Data),ok
	}
	return Data{},false
}
func (c *Cache) add(key string,value Data)  {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Add(key,value)
}
