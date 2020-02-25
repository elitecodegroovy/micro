package memory

import (
	"fmt"
	"github.com/micro/micro/mem"
	"runtime"
	"time"
)

var (
	done = make(chan bool)
)

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v M", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v M", bToMb(m.TotalAlloc))
	fmt.Printf("\tHeapSys = %v M", bToMb(m.HeapSys))
	fmt.Printf("\tSys = %v M", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func GetOSMem() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Sys)
}

func bToMb(b uint64) float64 {
	return float64(b) / float64(1024) / float64(1024)
}

func DoGoMemoryStatics() {
	ticker := time.NewTicker(time.Millisecond * 10)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case <-ticker.C:
			PrintMemUsage()
		}
	}
}

func GetFreeCache() uint64 {
	v, _ := mem.VirtualMemory()
	return v.Free
}

func Done() {
	done <- true
}
