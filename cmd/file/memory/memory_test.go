package memory

import (
	"fmt"
	"github.com/micro/micro/mem"
	"testing"
)

func TestDoStaticMemory(t *testing.T) {
	//go DoStaticMemory()
	//time.Sleep(3*time.Second)
	//Done()
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %vM, Free:%vM, UsedPercent:%f%%\n", float64(v.Total)/float64(1024)/float64(1024),
		float64(v.Free)/float64(1024)/float64(1024), v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)
	fmt.Println(1024 * 1024)
}
