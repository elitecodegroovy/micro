package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

func testPool() {
	drain := true
	var p sync.Pool
	const N = 100
loop:
	for try := 0; try < 3; try++ {
		if try == 1 {
			break
		}
		var fin, fin1 uint32
		for i := 0; i < N; i++ {
			v := new(string)
			runtime.SetFinalizer(v, func(vv *string) {
				atomic.AddUint32(&fin, 1)
			})
			p.Put(v)
		}
		if drain {
			for i := 0; i < N; i++ {
				p.Get()
			}
		}
		for i := 0; i < 5; i++ {
			runtime.GC()
			time.Sleep(time.Duration(i*100+10) * time.Millisecond)
			// 1 pointer can remain on stack or elsewhere
			if fin1 = atomic.LoadUint32(&fin); fin1 >= N-1 {
				log.Println(">>>>", fin1)
				continue loop
			}
		}
		log.Fatalf("only %v out of %v resources are finalized on try %v", fin1, N, try)
	}
}

func doPoolByteBuffer() {
	var pool = sync.Pool{
		// New creates an object when the pool has nothing available to return.
		// New must return an interface{} to make it flexible. You have to cast
		// your type after getting it.
		New: func() interface{} {
			// Pools often contain things like *bytes.Buffer, which are
			// temporary and re-usable.
			return &bytes.Buffer{}
		},
	}
	// When getting from a Pool, you need to cast
	s := pool.Get().(*bytes.Buffer)
	// We write to the object
	s.Write([]byte("dirty"))
	// Then put it back
	pool.Put(s)

	// Pools can return dirty results

	// Get 'another' buffer
	s = pool.Get().(*bytes.Buffer)
	// Write to it
	s.Write([]byte("append"))
	// At this point, if GC ran, this buffer *might* exist already, in
	// which case it will contain the bytes of the string "dirtyappend"
	fmt.Println(s)
	// So use pools wisely, and clean up after yourself
	s.Reset()
	pool.Put(s)

	// When you clean up, your buffer should be empty
	s = pool.Get().(*bytes.Buffer)
	// Defer your Puts to make sure you don't leak!
	defer pool.Put(s)
	s.Write([]byte("reset!"))
	// This prints "reset!", and not "dirtyappendreset!"
	fmt.Println(s)
}
