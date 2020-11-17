package main

import "sync"

type (
	Key   = string
	Value = string
)

type Cache interface {
	GetOrSet(key Key, valueFn func() Value) Value
	Get(key Key) (Value, bool)
}

type InMemoryCache struct {
	dataMutex sync.RWMutex
	data      map[Key]Value
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]
	return value, found
}

// GetOrSet возвращает значение ключа в случае его существования.
// Иначе, вычисляет значение ключа при помощи valueFn, сохраняет его в кэш и возвращает это значение.
func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {
	// TODO
}

func main() {

}
