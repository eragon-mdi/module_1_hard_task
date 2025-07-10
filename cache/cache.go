package cache

import "sync"

// Interface - реализуйте этот интерфейс
type Interface interface {
	Set(k, v string)
	Get(k string) (v string, ok bool)
}

// Не меняйте названия структуры и название метода создания экземпляра Cache, иначе не будут проходить тесты
type Cache struct {
	cache sync.Map
}

// NewCache создаёт и возвращает новый экземпляр Cache.
func NewCache() Interface {
	return &Cache{
		cache: sync.Map{},
	}
}

func (c *Cache) Set(k, v string) {
	c.cache.Store(k, v)
}

func (c *Cache) Get(k string) (v string, ok bool) {
	val, stored := c.cache.Load(k)
	if !stored {
		return "", false
	}

	v, ok = val.(string)
	return v, ok
}
