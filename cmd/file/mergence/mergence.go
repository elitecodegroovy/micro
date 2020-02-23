package mergence

import (
	"bufio"
	"bytes"
	"fmt"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
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
}

func New(outputFilePath string, lines int64, names []string) *mergence {
	return &mergence{
		targetFilePath: outputFilePath,
		filenames:      names,
		totalFileLines: lines,
		sortedDataChan: make(chan int64, lines*3/int64(len(names))/int64(100)),
	}
}

func (m *mergence) getChunkChanSize() int64 {
	return m.totalFileLines / int64(len(m.filenames)) / int64(100)
}

func (f *fileChunk) getFileBufferSize() int64 {
	return int64(os.Getpagesize() * 128)
}

func (f *fileChunk) readLinesFromFileBuffer() error {
	for {
		bytesLine, err := f.fileBuffer.ReadBytes('\n')
		if err == io.EOF {
			f.fileBuffer.Reset()
			break
		}
		if err != nil {
			fmt.Printf("Couldn't read bytes from file buffer of file %s : %v", f.filename, err)
			return err
		}

		e, err := strconv.ParseInt(strings.Trim(string(bytesLine), " \n"), 10, 64)
		if err != nil {
			fmt.Printf("Couldn't  parse bytesLine for the file %s : %v", f.filename, err)
			return err
		}
		f.chunkChan <- e
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

	bufBulk := make([]byte, f.getFileBufferSize())
	for {
		//Read bulk from file
		size, err := file.Read(bufBulk)
		if err == io.EOF {
			close(f.chunkChan)
			break
		}
		if err != nil {
			fmt.Printf("Couldn't read file chunk %s %v \n", f.filename, err)
			return err
		}
		f.fileBuffer = bytes.NewBuffer(bufBulk[:size])
		if err = f.readLinesFromFileBuffer(); err != nil {
			fmt.Printf("Couldn't read lines from file buffer bulk for the file %s %v \n", f.filename, err)
			return err
		}
	}
	return nil
}

func (m *mergence) Merge() error {
	var fileChunkSlice []*fileChunk
	var g errgroup.Group
	for i, filename := range m.filenames {
		fileChunkProcessing := &fileChunk{
			index:       i,
			filename:    filename,
			chunkChan:   make(chan int64, m.getChunkChanSize()),
			isProcessed: false,
		}
		fileChunkSlice = append(fileChunkSlice, fileChunkProcessing)

		//read file lines and send line data to the chan 'chunkChan'
		g.Go(func() error {
			return fileChunkProcessing.SendDataToFileChunkChan()
		})
	}
	m.fileChunkSlice = fileChunkSlice

	//get the sorted element and send to the chan 'sortedDataChan'
	g.Go(func() error {
		return m.compareLineValue()
	})

	//write data to the target file.
	g.Go(func() error {
		return m.writeFile()
	})

	if err := g.Wait(); err != nil {
		return err
	}
	return nil
}

func (m *mergence) compareLineValue() error {
	var comparingDataSlice []int64
	needAddedElements := 0
	var validDataChanMap = make(map[int]int)
	//The map can record the valid chan that can be receive data.
	for i, _ := range m.fileChunkSlice {
		validDataChanMap[i] = i
	}

	for i, chunkFileProcessingIndex := range m.fileChunkSlice {
		if !chunkFileProcessingIndex.isProcessed {
			e, ok := <-chunkFileProcessingIndex.chunkChan
			if !ok {
				chunkFileProcessingIndex.isProcessed = true
			}
			comparingDataSlice = append(comparingDataSlice, e)
		} else {
			delete(validDataChanMap, i)
			needAddedElements++
		}
	}
	if needAddedElements != 0 && len(validDataChanMap) != 0 {
	loopCallChan:
		for k, _ := range validDataChanMap {
			if needAddedElements != 0 {
				chunkFileProcessingIndex := m.fileChunkSlice[k]
				e, ok := <-chunkFileProcessingIndex.chunkChan
				if !ok {
					chunkFileProcessingIndex.isProcessed = true
					delete(validDataChanMap, k)
				} else {
					comparingDataSlice = append(comparingDataSlice, e)
					needAddedElements--
				}
			}
		}
		if needAddedElements != 0 && len(validDataChanMap) != 0 {
			goto loopCallChan
		}
	}

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
		if m.totalFileLines == 0 {
			break
		}
		e := <-m.sortedDataChan
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
	fmt.Printf("writing sorted elemented to the target file '%s' time: %f \n",
		m.targetFilePath,
		time.Since(t).Seconds())

	return nil
}
