package buf

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	writerBufferSize = getFileBufferSize() * 128
	ErrBufferRead    = errors.New("bytes.Buffer: couldn't read file chunk")
	ErrByteRead      = errors.New("couldn't read bytes from file buffer of file")
	ErrLineRead      = errors.New("couldn't  parse bytesLine for the file ")
)

// IntSlice attaches the methods of Interface to []int, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type FileChunkInfo struct {
	filename   string
	brokenLine []byte
	lines      int64
	lineItems  []int64
	fileBuffer *bytes.Buffer
}

func New(filename string) *FileChunkInfo {
	return &FileChunkInfo{
		filename:  filename,
		lineItems: []int64{},
	}
}

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

func getFileBufferSize() int64 {
	return int64(os.Getpagesize() * 128)
}

func DoSortOpts(filename string) string {
	var outputFileName string
	var err error
	sorter := New(filename)
	if outputFileName, err = sorter.ReadFileByBulkBuffer(); err != nil {
		fmt.Println("err: " + err.Error())
	}
	return outputFileName

}

func (f *FileChunkInfo) ReadFileByBulkBuffer() (string, error) {
	//t := time.Now()
	file, err := os.Open(f.filename)
	if err != nil {
		println(err)
	}
	defer file.Close()

	bufBulk := make([]byte, getFileBufferSize())
	for {
		//Read bulk from file
		size, err := file.Read(bufBulk)
		//fmt.Printf(">>>>>>>>>size:%d \n", size)
		if err == io.EOF {
			//fmt.Println("***file>" + f.filename + ", lines :" + strconv.FormatInt(f.lines, 10) + "\n")
			break
		}
		if err != nil {
			return "", ErrBufferRead
		}
		if len(f.brokenLine) > 0 {
			f.fileBuffer = bytes.NewBuffer(append(f.brokenLine, bufBulk[:size]...))
			f.brokenLine = []byte{}
		} else {
			f.fileBuffer = bytes.NewBuffer(bufBulk[:size])
		}

		for {
			bytesLine, err := f.fileBuffer.ReadBytes('\n')
			if err == io.EOF {
				f.brokenLine = bytesLine
				f.fileBuffer.Reset()
				break
			}
			if err != nil {
				return "", ErrByteRead
			}

			v := strings.Trim(string(bytesLine), " \n")
			if len(v) == 0 {
				continue
			}
			e, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				fmt.Println("------------v:[", v+"]")
				return "", ErrLineRead
			}
			f.lines++
			f.lineItems = append(f.lineItems, e)
		}
	}
	//fmt.Printf("read file time: %f \n", time.Since(t).Seconds())

	return sortOriginalData(f.lineItems, f.filename)
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

//writeDataToFile sort order time: 302.858123 for 1.2G
//writing sort order time: 3.431778  for 15M
func writeDataToFile(int64DataSlice []int64, filename string) (string, error) {
	//t := time.Now()
	baseDir, inputFilename, ext := GetFileNameInfo(filename)
	var sortedFilenamePath string
	if len(baseDir) == 0 {
		sortedFilenamePath = inputFilename + "Sorted." + ext
	} else {
		sortedFilenamePath = baseDir + string(filepath.Separator) + inputFilename + "Sorted." + ext
	}

	f, err := os.OpenFile(sortedFilenamePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		fmt.Printf("failed to open file : %s ", err.Error())
		return "", err
	}

	//100K buffer writer size,
	w := bufio.NewWriterSize(f, int(writerBufferSize))
	//w := bufio.NewWriter(f)
	for _, item := range int64DataSlice {
		//print(">")
		if _, err1 := w.Write([]byte(strconv.FormatInt(item, 10) + "\n")); err1 != nil {
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
	//fmt.Printf("writing sort order time: %f \n", time.Since(t).Seconds())

	//clean the buff/cache item in OS
	int64DataSlice = []int64{}
	//memory.CleanBufferCacheOfOS()
	return sortedFilenamePath, nil
}

//GetFileNameInfo extracts name and extension from path
func GetFileNameInfo(path string) (string, string, string) {
	split := strings.Split(path, string(filepath.Separator))
	var fileBase, name string
	if len(split) == 1 {
		fileBase = ""
	} else if len(split) == 2 {
		fileBase = split[0]
	} else {
		//file directory path
		fileBase = strings.Join(split[:len(split)-2], string(filepath.Separator))
	}
	//file name
	name = split[len(split)-1]

	split = strings.Split(name, ".")
	ext := split[len(split)-1]
	nameSlice := split[:len(split)-1]
	name = strings.Join(nameSlice, "")

	return fileBase, name, ext
}
