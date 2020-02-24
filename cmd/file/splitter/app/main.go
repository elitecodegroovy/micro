package main

import (
	"flag"
	"fmt"
	"github.com/micro/micro/cmd/file/splitter"
	"os"
)

func main() {
	var inputFilenamePath = flag.String("filenamePath", "bigLongTypeData.txt", "provide the filename for splitter")
	var outputFileDir = flag.String("outputDir", "data/default", "provide the filename for splitter")
	var fileChunkSize = flag.Int64("fileChunkSize", 10485760, "provide the file chunk size for splitter (B unit), default 10485760")
	flag.Parse()

	splitter := splitter.New()
	splitter.FileChunkSize = *fileChunkSize //10M
	result, _ := splitter.Split(*inputFilenamePath, *outputFileDir)
	if len(result) == 0 {
		fmt.Printf("failed ...\n")
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
		fmt.Printf("%s, size: %fM", item, float64(statChunkFile.Size())/float64(1024)/float64(1024))
	}

}
