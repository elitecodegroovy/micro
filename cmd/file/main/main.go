package main

import (
	"flag"
	"fmt"
	buf "github.com/micro/micro/cmd/file/buf"
	"github.com/micro/micro/cmd/file/memory"
	"github.com/micro/micro/cmd/file/mergence"
	"github.com/micro/micro/cmd/file/splitter"
	"log"
	"os"
	"time"
)

var filename = flag.String("filename", "data/bigLongTypeData.txt", "provide the input filename, default filename 'data/bigLongTypeData.txt'")
var maxChunkFileSize = flag.Int64("maxChunkFileSize", int64(memory.GetFreeCache()/3), "max chunk file size, default the `osFreeCache/3`")
var sortedFilename = flag.String("sorted-filename", "data/sortedBigLongTypeData.txt", "provide the output filename, default sorted data file filename 'data/sortedBigLongTypeData.txt'")
var tempOutputDir = flag.String("temp-output-dir", "data/default", "provide the ouput data directory , default is directory 'data/default'")
var debug = flag.Bool("debug", false, "print the debug info, default is false.")

//validateInputParameters does the validation work in order to make sure that free memory is
//enough for file I/O opts.
func validateInputParameters() bool {
	if *maxChunkFileSize > int64(memory.GetFreeCache()/2) {
		fmt.Println("each chunk file size is too big. please provide a valid parameter 'split-n', which default is 10.")
		return false
	}
	return true
}

func splitFile() ([]string, int64, error) {
	splitter := splitter.New()

	splitter.FileChunkSize = *maxChunkFileSize
	fmt.Println("os RAM free ", memory.GetFreeCache()/1024/1024, "M, max chunk file size :", splitter.FileChunkSize/1024/1024, "M")

	result, err := splitter.Split(*filename, *tempOutputDir)
	if err != nil {
		fmt.Println("do split big file with an error :", err.Error())
		return nil, 0, err
	}
	return result, splitter.BigFileTotalLineNumber, nil
}

func clearTempFiles(filenames []string) {
	if len(filenames) > 0 {
		for _, filename := range filenames {
			if err := os.Remove(filename); err != nil {
				fmt.Println("os.Remove file path ", filename, " with an error. ", err.Error())
			}
		}
	}
}

func main() {
	t := time.Now()
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
	//fmt.Println("output file info :", result)
	if len(result) == 0 {
		fmt.Println("failed ...error occurred in server internal.", result)
	}

	var sortedFilenamePaths []string
	//sort the data in each small file
	for _, resultFilepath := range result {

		sorter := buf.New(resultFilepath)
		sortedFilenamePath, err := sorter.ReadFileByBulkBuffer()
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
	fmt.Printf(">>>total lines(include the empty line): %d \n", fileLines)
	if len(sortedFilenamePaths) != 0 {
		m := mergence.New(*sortedFilename, fileLines, sortedFilenamePaths)
		if err = m.Merge(); err != nil {
			log.Fatal("..." + err.Error())
		}
	}
	fmt.Println("--------------------------------------------")
	fmt.Printf("----------total time elapsed %fs------------\n", time.Since(t).Seconds())
	fmt.Println("--------------------------------------------")

	clearTempFiles(append(result, sortedFilenamePaths...))
	//memory.CleanBufferCacheOfOS()
}
