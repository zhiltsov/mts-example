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
/*
TODO Это вариант №1 без исправления кода выше
TODO Минусом логики является блокирующя операция на переданное замыкание. Они выполняются последовательно
TODO Необходимо выполнение замыканий выделять в отдельный слой управления
*/
func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {
	if value, found := cache.Get(key); found { // если ключ в кеше - отдаем
		return value
	}

	cache.dataMutex.Lock() // блокируем весь кеш
	defer cache.dataMutex.Unlock()
	if value, found := cache.data[key]; found { // еще раз проверяем ключ после блокировки всех горутин
		return value // отдаем, если одна из горутин записала ключ, чтобы не выполнять функцию повторно
	}
	cache.data[key] = valueFn()
	return cache.data[key]
}
