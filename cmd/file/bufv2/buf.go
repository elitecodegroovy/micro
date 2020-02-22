package buf

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/sync/errgroup"
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

	bufferSize = int64(os.Getpagesize() * 128)
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
	//set buffer block size 10M
	buf := make([]byte, 0, 1024*1024*10)

	// Buffer sets the initial buffer to use when scanning and the maximum
	// size of buffer that may be allocated during scanning.
	fmt.Println("local memory : ", getOSMem())
	//1G cap
	scanner.Buffer(buf, 1024*1024*1024)

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

	//sortOriginalData(originalData, filename)

}

func sortOriginalData(originalData []int64, filename string) {
	t := time.Now()
	//quick sort algorithm
	int64Slice := Int64Slice(originalData)
	sort.Sort(int64Slice)
	//selectionSort(originalData)

	fmt.Printf("sort time: %f \n", time.Since(t).Seconds())
	//ShowMemoryInfo()
	writeDataToFile(originalData, "sorted"+filename)
	//Done()
}

//writing sort order time: 302.858123 for 1.2G
//writing sort order time: 3.431778  for 15M
func writeDataToFile(int64DataSlice []int64, filename string) {
	t := time.Now()
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, os.ModePerm)

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

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %f M", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %fM", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %f M", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %d \n", m.NumGC)
}

func getOSMem() int {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return int(m.Sys)
}

func bToMb(b uint64) float64 {
	return float64(b) / float64(1024) / float64(1024)
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

func Done() {
	done <- true
}

func ShowMemoryInfo() {
	go doStaticMemory()
}
