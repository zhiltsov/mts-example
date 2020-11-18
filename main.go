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

type ChannelRunner struct {
	chMutex  sync.Mutex
	chFunc   chan func() Value
	chResult chan Value
}

type InMemoryCache struct {
	dataMutex sync.RWMutex
	data      map[Key]Value
	ch        map[Key]*ChannelRunner
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[Key]Value),
		ch:   make(map[Key]*ChannelRunner),
	}
}

func (cache *InMemoryCache) Get(key Key) (Value, bool) {
	cache.dataMutex.RLock()
	defer cache.dataMutex.RUnlock()

	value, found := cache.data[key]
	return value, found
}

func (cache *InMemoryCache) Set(key Key, value Value) {
	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()

	cache.data[key] = value
}

func (cache *InMemoryCache) MakeRunner(key Key) *ChannelRunner {
	cache.dataMutex.Lock()
	defer cache.dataMutex.Unlock()

	if _, found := cache.ch[key]; !found { // Используем для каждого ключа свой канал
		cache.ch[key] = &ChannelRunner{ // Буфера нет!
			chFunc:   make(chan func() Value),
			chResult: make(chan Value),
		}
	}
	return cache.ch[key]
}

// GetOrSet возвращает значение ключа в случае его существования.
// Иначе, вычисляет значение ключа при помощи valueFn, сохраняет его в кэш и возвращает это значение.
/*
Вариант №2 c исправлением кода выше исходного для реализации асинхронности замыканий
*/
func (cache *InMemoryCache) GetOrSet(key Key, valueFn func() Value) Value {
	runner := cache.MakeRunner(key)

	go func() { // обработчик данных кеша в горутине
		runner.chResult <- func(fn func() Value) (value Value) {
			if value, found := cache.Get(key); found { // если ключ в кеше - отдаем
				return value
			}
			value = fn()
			cache.Set(key, value) // установка нового значения в кеш
			return
		}(<-runner.chFunc)
	}()

	runner.chMutex.Lock()
	defer runner.chMutex.Unlock()

	runner.chFunc <- valueFn
	return <-runner.chResult
}
