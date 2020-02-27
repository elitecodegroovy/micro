package mergence

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type mergence struct {
	filenames        []string
	totalFileLines   int64
	writerBufferSize int
	sortedDataChan   chan int64
	targetFilePath   string
	fileChunkSlice   []*fileChunk
}

type fileChunk struct {
	index       int
	filename    string
	chunkChan   chan int64
	fileBuffer  *bytes.Buffer
	isProcessed bool
	lines       int64
	brokenLine  []byte
}

type lineItem struct {
	index int
	value int64
}

type lineItemSlice []lineItem

func (p lineItemSlice) Len() int           { return len(p) }
func (p lineItemSlice) Less(i, j int) bool { return p[i].value < p[j].value }
func (p lineItemSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func New(outputFilePath string, lines int64, names []string) *mergence {
	return &mergence{
		targetFilePath: outputFilePath,
		filenames:      names,
		totalFileLines: lines,
		sortedDataChan: make(chan int64, int(lines*2)/len(names)/100),
	}
}

func (m *mergence) getChunkChanSize() int64 {
	return m.totalFileLines / int64(len(m.filenames)) / int64(100)
}

func getFileBufferSize() int64 {
	return int64(os.Getpagesize() * 128)
}

func (f *fileChunk) readLinesFromFileBuffer() error {
	for {
		bytesLine, err := f.fileBuffer.ReadBytes('\n')
		if err == io.EOF {
			f.brokenLine = bytesLine
			f.fileBuffer.Reset()
			break
		}
		if err != nil && err != io.EOF {
			fmt.Printf("Couldn't read bytes from file buffer of file %s : %v", f.filename, err)
			return err
		}

		v := strings.Trim(string(bytesLine), " \n")
		e, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			fmt.Printf("Couldn't  parse bytesLine for the file %s : %v. \n %s", f.filename, err, v)
			return err
		}
		f.chunkChan <- e
		f.lines++
		//fmt.Printf("i: %d, e:%d\n", f.lines, e)
	}
	return nil
}

func (f *fileChunk) SendDataToFileChunkChan() error {
	file, err := os.Open(f.filename)
	if err != nil {
		fmt.Printf("Couldn't open file  %s %v \n", f.filename, err)
		return err
	}
	defer file.Close()

	bufBulk := make([]byte, getFileBufferSize())
	for {
		//Read bulk from file
		size, err := file.Read(bufBulk)
		if err == io.EOF {
			close(f.chunkChan)
			//fmt.Println("***close chan >" + f.filename + ", lines :" + strconv.FormatInt(f.lines, 10))
			break
		}
		if err != nil && err != io.EOF {
			fmt.Printf("Couldn't read file chunk %s %v \n", f.filename, err)
			return err
		}
		if len(f.brokenLine) > 0 {
			f.fileBuffer = bytes.NewBuffer(append(f.brokenLine, bufBulk[:size]...))
			f.brokenLine = []byte{}
		} else {
			f.fileBuffer = bytes.NewBuffer(bufBulk[:size])
		}

		if err = f.readLinesFromFileBuffer(); err != nil {
			fmt.Printf("Couldn't read lines from file buffer bulk for the file %s %v \n", f.filename, err)
			return err
		}
	}
	return nil
}

//Merge merges all files into a sorted data file.
func (m *mergence) Merge() error {
	m.fileChunkSlice = []*fileChunk{}
	var g errgroup.Group
	for i, filename := range m.filenames {
		fileChunkProcessing := &fileChunk{
			index:       i,
			filename:    filename,
			chunkChan:   make(chan int64, m.getChunkChanSize()),
			isProcessed: false,
		}
		m.fileChunkSlice = append(m.fileChunkSlice, fileChunkProcessing)

		//Step 1:
		//read file lines and send line data to the chan 'chunkChan'
		g.Go(func() error {
			return fileChunkProcessing.SendDataToFileChunkChan()
		})

		//g.Go(func()error {
		//	e, ok := <- fileChunkProcessing.chunkChan
		//	for ok {
		//		fmt.Printf(">>> i :%d, e: %d", fileChunkProcessing.lines, e)
		//		e, ok = <- fileChunkProcessing.chunkChan
		//	}
		//	return nil
		//})
	}

	//Step  2:
	//get the sorted element and send to the chan 'sortedDataChan'
	g.Go(func() error {
		return m.compareLineValueForEachFile()
	})

	//Step  3:
	//write data to the target file.
	g.Go(func() error {
		return m.writeFile()
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func minInSlice(a []int64) int {
	if a == nil {
		return -1
	}
	if len(a) == 1 {
		return 0
	}
	i := 1
	min := 0
	for i < len(a) {
		if a[min] > a[i] {
			min = i
		}
		i++
	}
	return min
}

func sortItems(items []lineItem) {
	if len(items) == 1 {
		return
	}
	s := lineItemSlice(items)
	sort.Sort(s)
}

func (m *mergence) compareLineValueForEachFile() error {
	var fileFirstLineValueSlice []lineItem
	//The map can record the valid chan that can be receive data.
	for i, chunkFileProcessingIndex := range m.fileChunkSlice {
		e, ok := <-chunkFileProcessingIndex.chunkChan
		if ok {
			item := lineItem{
				index: i,
				value: e,
			}
			fileFirstLineValueSlice = append(fileFirstLineValueSlice, item)
		}
	}
loop:
	if len(fileFirstLineValueSlice) > 1 {
		sortItems(fileFirstLineValueSlice)

		//send item value 'e'
		m.sortedDataChan <- fileFirstLineValueSlice[0].value
		e, ok := <-m.fileChunkSlice[fileFirstLineValueSlice[0].index].chunkChan
		if !ok {
			m.fileChunkSlice[fileFirstLineValueSlice[0].index].isProcessed = true
			//remove the first line
			fileFirstLineValueSlice = append(fileFirstLineValueSlice[:0], fileFirstLineValueSlice[1:]...)
		} else {
			fileFirstLineValueSlice[0].value = e
		}
		goto loop
	} else if len(fileFirstLineValueSlice) == 1 {
		m.sortedDataChan <- fileFirstLineValueSlice[0].value
		e, ok := <-m.fileChunkSlice[fileFirstLineValueSlice[0].index].chunkChan
		for ok {
			//fmt.Println(">>>>>>", e)
			//fmt.Printf("+++: %d\n", e)
			m.sortedDataChan <- e
			e, ok = <-m.fileChunkSlice[fileFirstLineValueSlice[0].index].chunkChan
		}

		close(m.sortedDataChan)
	}
	//fmt.Printf(">>>>*****len(fileFirstLineValueSlice): %d \n", len(fileFirstLineValueSlice))

	return nil
}

func (m *mergence) writeFile() error {
	t := time.Now()
	f, err := os.OpenFile(m.targetFilePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		fmt.Printf("failed to open file : %s ", err.Error())
		return err
	}

	w := bufio.NewWriterSize(f, m.writerBufferSize)
	for {
		e, ok := <-m.sortedDataChan
		if !ok {
			//fmt.Printf("<<<<<<>>>>>>close chan 'sortedDataChan': %d\n", m.totalFileLines)
			break
		}
		//fmt.Printf("------: %d, e: %d\n", m.totalFileLines, e)
		m.totalFileLines--
		if _, err1 := w.Write([]byte(strconv.FormatInt(e, 10) + "\n")); err1 != nil {
			fmt.Printf("failed to write number data : %s \n", err1.Error())
			return err
		}
	}

	if err = w.Flush(); err != nil {
		fmt.Printf("w.Flush() error  : %s \n ", err.Error())
		return err
	}
	if err = f.Close(); err != nil {
		fmt.Printf("w.Flush() error  : %s \n ", err.Error())
		return err
	}
	fmt.Printf("\noutput file path : '%s' , writing time: %f \n",
		m.targetFilePath,
		time.Since(t).Seconds())

	return nil
}
