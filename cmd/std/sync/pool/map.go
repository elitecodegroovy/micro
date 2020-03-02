package main

import (
	"log"
	"math/rand"
	"runtime"
	"sync"
)

func syncMap() {
	const mapSize = 1 << 10

	m := new(sync.Map)
	for n := int64(1); n <= mapSize; n++ {
		m.Store(n, int64(n))
	}

	done := make(chan struct{})
	var wg sync.WaitGroup
	defer func() {
		close(done)
		wg.Wait()
	}()
	//Default cpu core number . the value 0 you can set for 'GOMAXPROCS'.
	for g := int64(runtime.GOMAXPROCS(0)); g > 0; g-- {
		r := rand.New(rand.NewSource(g))
		wg.Add(1)
		go func(g int64) {
			defer wg.Done()
			for i := int64(0); ; i++ {
				select {
				case <-done:
					return
				default:
				}
				for n := int64(1); n < mapSize; n++ {
					if r.Int63n(mapSize) == 0 {
						m.Store(n, n*i*g)
					} else {
						m.Load(n)
					}
				}
			}
		}(g)
	}

	iters := 1 << 10
	for n := iters; n > 0; n-- {
		seen := make(map[int64]bool, mapSize)

		m.Range(func(ki, vi interface{}) bool {
			k, v := ki.(int64), vi.(int64)
			if v%k != 0 {
				log.Fatalf("while Storing multiples of %v, Range saw value %v", k, v)
			}
			if seen[k] {
				log.Fatalf("Range visited key %v twice", k)
			}
			seen[k] = true

			return true
		})
		if len(seen) != mapSize {
			log.Fatalf("Range visited %v elements of %v-element Map", len(seen), mapSize)
		}
	}

	mapRange()
}

func mapRange() {
	m := new(sync.Map)

	m.Store(1, 0)
	m.Store(2, 1)
	m.Store(3, 2)
	m.Range(func(ki, vi interface{}) bool {
		k, v := ki.(int), vi.(int)
		log.Printf("k: %d, v: %d", k, v)

		return true
	})
}
