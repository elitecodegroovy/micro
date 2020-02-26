package main

import (
	"flag"
	"fmt"
	"github.com/micro/micro/cmd/file/memory"
	"github.com/micro/micro/cmd/file/splitter"
	"os"
)

func main() {
	var inputFilenamePath = flag.String("filenamePath", "bigLongTypeData.txt", "provide the filename for splitter")
	var outputFileDir = flag.String("outputDir", "data/default", "provide the filename for splitter")
	var fileChunkSize = flag.Int64("fileChunkSize", int64(memory.GetFreeCache()/2), "provide the file chunk size for splitter (B unit), default 10485760")
	flag.Parse()

	splitter := splitter.New()
	splitter.FileChunkSize = *fileChunkSize //10M
	result, err := splitter.Split(*inputFilenamePath, *outputFileDir)
	if err != nil {
		fmt.Printf("failed ...%s", err.Error())
		os.Exit(1)
	}

	for _, item := range result {
		chunkFile, err := os.Open(item)
		if err != nil {
			fmt.Printf("err: %s \n", err.Error())
		}

		statChunkFile, err := chunkFile.Stat()
		if err != nil {
			fmt.Printf("err: %s \n", err.Error())
		}
		fmt.Printf("%s, size: %fM\n", item, float64(statChunkFile.Size())/float64(1024)/float64(1024))
	}

}
