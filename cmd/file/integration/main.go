package main

import (
	"flag"
	"fmt"
	buf "github.com/micro/micro/cmd/file/bufv2"
	"github.com/micro/micro/cmd/file/memory"
	"github.com/micro/micro/cmd/file/splitter"
	"os"
)

var filename = flag.String("filename", "data/bigLongTypeData.txt", "provide the input filename, default filename 'data/bigLongTypeData.txt'")
var splitN = flag.Int64("split-n", 5, "split to N file.")
var sortedFilename = flag.String("sorted-filename", "data/sortedBigLongTypeData.txt", "provide the output filename, default sorted data file filename 'data/sortedBigLongTypeData.txt'")
var outputDir = flag.String("output-dir", "data/default", "provide the ouput data directory , default is directory 'data/default'")
var chunkFileSize int64

//validateInputParameters does the validation work in order to make sure that free memory is
//enough for file I/O opts.
func validateInputParameters() bool {
	//max allowed file size.
	var maxChunkFileSize int64
	file, err := os.Open(*filename)
	if err != nil {
		fmt.Println(" os.Open filename error :", err.Error())
		return false
	}
	defer file.Close()

	fileInfo, err := os.Stat(*filename)
	if err != nil {
		fmt.Println("os.Stat filename error :", err.Error())
		return false
	}
	fileSize := fileInfo.Size()

	if fileSize%(*splitN) == 0 {
		maxChunkFileSize = fileSize / (*splitN)

	} else {
		addedPartition := fileSize % (*splitN)
		maxChunkFileSize = addedPartition + fileSize/(*splitN)
	}
	if maxChunkFileSize > int64(memory.GetFreeCache()/2) {
		fmt.Println("each chunk file size is too big. please provide a valid parameter 'split-n', which default is 10.")
		return false
	}
	chunkFileSize = maxChunkFileSize

	return true
}

func splitFile() ([]string, int64, error) {
	splitter := splitter.New()

	splitter.FileChunkSize = chunkFileSize
	fmt.Println("max chunk file size :", splitter.FileChunkSize, "K")

	result, err := splitter.Split(*filename, *outputDir)
	if err != nil {
		fmt.Println("do split big file with an error :", err.Error())
		return nil, 0, err
	}
	return result, splitter.BigFileTotalLineNumber, nil
}

func main() {
	flag.Parse()

	if !validateInputParameters() {
		fmt.Println("please input cmd parameters --filename --split-n --sorted-filename --output-dir")
		return
	}

	//Step 1: split big file into small file
	result, fileLines, err := splitFile()
	if err != nil {
		fmt.Println("splitFile splits big file with an error :", err.Error())
		return
	}
	fmt.Println("output file info :", result)
	if len(result) == 0 {
		fmt.Println("failed ...error occurred in server internal.", result)
	}

	var sortedFilenamePaths []string
	//sort the data in each small file
	for i, resultFilepath := range result {
		fmt.Printf("%d > sort file %s \n", i, resultFilepath)
		sortedFilenamePath, err := buf.SortDataInFile(resultFilepath, int(chunkFileSize/3), int(chunkFileSize))
		if err != nil {
			fmt.Println("buf.SortDataInFile :", err.Error())
			return
		}
		if len(sortedFilenamePath) != 0 {
			sortedFilenamePaths = append(sortedFilenamePaths, sortedFilenamePath)
		} else {
			fmt.Println("it failed to get sorted file name:")
			return
		}
	}
	fmt.Printf(">>>total lines: %d \n", fileLines)

}
