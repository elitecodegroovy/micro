package buf

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/sync/errgroup"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	bufferSize       = int64(os.Getpagesize() * 128)
	writerBufferSize = 102400
)

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// It is a good choice if the elements is less than 12.
func selectionSort(items []int64) {
	var n = len(items)
	for i := 0; i < n; i++ {
		var minIdx = i
		for j := i; j < n; j++ {
			if items[j] < items[minIdx] {
				minIdx = j
			}
		}
		items[i], items[minIdx] = items[minIdx], items[i]
	}
}

func ReadLinesByBufIO(filename string, bufferSize, bufferSizeCapacity int) (string, error) {
	writerBufferSize = bufferSize

	t := time.Now()
	file, err := os.Open(filename)
	if err != nil {
		println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//set buffer block size 10M
	buf := make([]byte, 0, bufferSize)

	// Buffer sets the initial buffer to use when scanning and the maximum
	// size of buffer that may be allocated during scanning.
	//1G cap
	scanner.Buffer(buf, bufferSizeCapacity)

	var i int64
	//16 M
	var originalData []int64
	for scanner.Scan() {
		i++
		line := strings.Trim(scanner.Text(), " \n")
		if line == "" {
			continue
		}
		//fmt.Printf("%d %s", i, scanner.Text())

		number, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Printf("[%s]\n", line)
		}
		originalData = append(originalData, number)
	}
	fmt.Printf("read file time: %f \n", time.Since(t).Seconds())

	return sortOriginalData(originalData, filename)
}

func sortOriginalData(originalData []int64, filename string) (string, error) {
	t := time.Now()
	//quick sort algorithm
	int64Slice := Int64Slice(originalData)
	sort.Sort(int64Slice)
	//selectionSort(originalData)

	fmt.Printf("sort time: %f \n", time.Since(t).Seconds())
	//ShowMemoryInfo()
	return writeDataToFile(originalData, filename)
	//Done()
}

//writing sort order time: 302.858123 for 1.2G
//writing sort order time: 3.431778  for 15M
func writeDataToFile(int64DataSlice []int64, filename string) (string, error) {
	t := time.Now()
	baseBase, inputFilename, ext := getFileNameInfo(filename)
	sortedFilenamePath := baseBase + string(filepath.Separator) + inputFilename + "Sorted." + ext
	f, err := os.OpenFile(sortedFilenamePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		fmt.Printf("failed to open file : %s ", err.Error())
		return "", err
	}

	//100K buffer writer size,
	//3.223659s
	w := bufio.NewWriterSize(f, writerBufferSize)
	//3.372287 s
	//w := bufio.NewWriter(f)
	for _, data := range int64DataSlice {
		//print(">")
		if _, err1 := w.Write([]byte(strconv.FormatInt(data, 10) + "\n")); err1 != nil {
			fmt.Printf("failed to write number data : %s \n", err1.Error())
			return "", err
		}
	}
	if err = w.Flush(); err != nil {
		fmt.Printf("w.Flush() error  : %s \n ", err.Error())
		return "", err
	}
	if err = f.Close(); err != nil {
		fmt.Printf("w.Flush() error  : %s \n ", err.Error())
		return "", err
	}
	fmt.Printf("writing sort order time: %f \n", time.Since(t).Seconds())
	return sortedFilenamePath, nil
}

//It extracts name and extension from path
func getFileNameInfo(path string) (string, string, string) {
	split := strings.Split(path, string(filepath.Separator))
	//file directory path
	fileBase := strings.Join(split[:len(split)-2], string(filepath.Separator))
	//file name
	name := split[len(split)-1]

	split = strings.Split(name, ".")
	ext := split[len(split)-1]
	nameSlice := split[:len(split)-1]
	name = strings.Join(nameSlice, "")

	return fileBase, name, ext
}

// 314.925621
func writeDataToFileV0(int64DataSlice []int64, filename string) {
	t := time.Now()
	var g errgroup.Group
	bulkBuffer := bytes.NewBuffer(make([]byte, 0, bufferSize))
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		fmt.Printf("failed to open file : %s ", err.Error())
		os.Exit(1)
	}

	orderDataChan := make(chan int64, bufferSize)
	defer close(orderDataChan)

	sliceSize := len(int64DataSlice)
	g.Go(func() error {
		for len(int64DataSlice) > 0 {
			orderDataChan <- int64DataSlice[0]
			int64DataSlice = append(int64DataSlice[:0], int64DataSlice[1:]...)
		}
		return nil
	})

	g.Go(func() error {
		for {
			if sliceSize == 0 {
				return nil
			}
			v, ok := <-orderDataChan
			if ok {
				sliceSize -= 1
				if _, err := bulkBuffer.Write([]byte(strconv.FormatInt(v, 10) + "\n")); err != nil {
					fmt.Printf("failed to write number data : %s ", err.Error())
					os.Exit(1)
				}
				if int64(bulkBuffer.Len()) >= (bufferSize - 1024) {
					if _, err = f.Write(bulkBuffer.Bytes()); err != nil {
						fmt.Printf("failed to write []byte for a file. Err : %s ", err.Error())
						os.Exit(1)
					}
					bulkBuffer.Reset()
				}
			}
		}
	})
	if err := g.Wait(); err != nil {
		fmt.Println("error ing.Wait() , error " + err.Error())
	}
	//writing sort order time: 389.195321
	fmt.Printf("writing sort order time: %f \n", time.Since(t).Seconds())
}
