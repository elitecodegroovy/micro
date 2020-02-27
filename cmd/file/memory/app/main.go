package main

import (
	"fmt"
	"github.com/micro/micro/cmd/file/memory"
)

func main() {
	fmt.Println("Run before OS Free Cache ", memory.GetFreeCache()/1024/1024, "M")
	memory.CleanBufferCacheOfOS()
	fmt.Println("Run after OS Free Cache ", memory.GetFreeCache()/1024/1024, "M")
}
