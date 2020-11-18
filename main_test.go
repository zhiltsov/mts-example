package main

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

const (
	GoCounts  = 1_000 // Количество горутин
	CacheSize = 10     // Количество ключей результирующего кеша
)

/*
Инкрементный счетчик записи в кеш
*/
type Counter struct {
	sync.Mutex
	v uint64
}

// +1
func (c *Counter) Inc() {
	c.Lock()
	defer c.Unlock()
	c.v++
}

var (
	cache   = NewInMemoryCache()
	counter = new(Counter)
)

/*
Тест: отдай или прими!
*/
func TestInMemoryCache_GetOrSet(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(GoCounts)

	fmt.Printf(
		"Запускаем %d горутин(у,ы) [GoCounts Х CacheSize] для попытки записи в кеш по %d значений(я,ю)\n",
		GoCounts*CacheSize,
		CacheSize,
	)
	fmt.Println("Время работы каждого замыкания - 1 секунда. time.Sleep(time.Second)")
	fmt.Printf("Ожидайте несколько секунд\n\n")

	for i := 0; i < GoCounts; i++ {
		i := i
		go func() { // Запускаем GoCounts горутин
			defer wg.Done()
			var fwg sync.WaitGroup
			fwg.Add(CacheSize)
			for k := 0; k < CacheSize; k++ { // Читаем или пишем CacheSize значений
				go func(i, k int) {
					defer fwg.Done()
					cache.GetOrSet(
						strconv.Itoa(k), // Ключ кеша
						func() Value { // Записываем в кеш, если ключ не найден
							fmt.Printf("Горутина %d, ключ: %d\n", i, k)
							time.Sleep(time.Second)    // Цена каждой записи - секунда
							counter.Inc()              // +1 к кол-ву записей в кеш
							return strconv.Itoa(k * k) // в результат квадрат
						})
				}(i, k)
			}
			fwg.Wait()
		}()
	}

	wg.Wait()
	if counter.v != CacheSize {
		t.Errorf("Кол-во обращений на запись (%d) в кеш не совпадает с размером (%d)", counter.v, CacheSize)
	} else {
		fmt.Println("OK! Кол-во записи в кеш равно его размеру.")
	}
}
