package splitter

import (
	"fmt"
)

func ExampleSplitFS() {
	splitter := New()
	splitter.FileChunkSize = 10485760 //10M
	result, _ := splitter.Split("/home/app/goapp/src/github.com/micro/micro/bigLongTypeData.txt", "/home/app/goapp/src/github.com/micro/micro/")
	fmt.Println(result)

}
