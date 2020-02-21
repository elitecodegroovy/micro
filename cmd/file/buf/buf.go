package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	fileName       = "bigLongTypeData.txt"
	outputFileName = "sortedDataSet.txt"
	done           = make(chan bool)
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

func ReadLinesByBufIO(filename string) {
	t := time.Now()
	file, err := os.Open(filename)
	if err != nil {
		println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	//set buffer block size 1M
	buf := make([]byte, 0, 1024*1024)

	// Buffer sets the initial buffer to use when scanning and the maximum
	// size of buffer that may be allocated during scanning.
	scanner.Buffer(buf, 130*1024*1024)

	var i int64
	//16 M
	var originalData []int64
	for scanner.Scan() {
		//if i > 1024*1024 {
		//	fmt.Printf("the line number is too large (bigger than 16*1024*1024)")
		//}
		i++
		line := strings.Trim(scanner.Text(), " \n")
		if line == "" {
			continue
		}
		//fmt.Printf("%d %s", i, line+ "\n")

		number, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Printf("[%s]\n", line)
		}
		originalData = append(originalData, number)

	}
	fmt.Printf("read file time: %f \n", time.Since(t).Seconds())
	t = time.Now()
	//quick sort algorithm
	int64Slice := Int64Slice(originalData)
	sort.Sort(int64Slice)
	//selectionSort(originalData)

	fmt.Printf("sort time: %f \n", time.Since(t).Seconds())
	writeDataToFile(originalData)
}

func writeDataToFile(int64DataSlice []int64) {
	t := time.Now()
	f, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		fmt.Printf("failed to open file : %s ", err.Error())
		os.Exit(1)
	}

	defer f.Close()
	//100K buffer writer size,
	//3.223659s
	w := bufio.NewWriterSize(f, 1024*100)
	//3.372287 s
	//w := bufio.NewWriter(f)
	for _, data := range int64DataSlice {
		//print(">")
		if _, err1 := f.Write([]byte(strconv.FormatInt(data, 10) + "\n")); err1 != nil {
			fmt.Printf("failed to write number data : %s ", err1.Error())
			os.Exit(1)
		}
	}
	w.Flush()
	fmt.Printf("writing sort order time: %f \n", time.Since(t).Seconds())
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v KB", bToKb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v KB", bToKb(m.TotalAlloc))
	fmt.Printf("\tSys = %v KB", bToKb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToKb(b uint64) float64 {
	return float64(b) / float64(1024)
}

func doStaticMemory() {
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

func showMemoryInfo() {
	go doStaticMemory()
}
func main() {
	showMemoryInfo()
	ReadLinesByBufIO(fileName)

	done <- true
}
