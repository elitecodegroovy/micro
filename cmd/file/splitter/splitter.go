//Package split_csv implements splitting of csv files on chunks by size in bytes
package splitter

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//minFileChunkSize min file chunk size in bytes
const minFileChunkSize = 100

var (
	ErrSmallFileChunkSize = errors.New("file chunk size is too small")
	ErrBigFileChunkSize   = errors.New("file chunk size is bigger than input file")
)

//Splitter struct which contains options for splitting
//FileChunkSize - a size of chunk in bytes, should be set by client
//WithHeader - whether split csv with header (true by default)
type Splitter struct {
	FileChunkSize int64 //in bytes
	WithHeader    bool
	bufferSize    int64 //in bytes
}

type ChunkFileInfo struct {
	s             Splitter
	inputFilePath string
	fileName      string
	ext           string
	outputDirPath string
	inputFile     *os.File
	chunkFile     *os.File
	chunkFilePath string
	header        []byte
	firstLine     bool
	brokenLine    []byte
	chunk         int
	bulkBuffer    *bytes.Buffer
	fileBuffer    *bytes.Buffer
	result        []string
}

func (s *ChunkFileInfo) setChunkFile(file *os.File) {
	if s.chunkFile != nil {
		s.chunkFile.Close()
	}
	s.chunkFile = file
}

//New initializes Splitter struct
func New() Splitter {
	return Splitter{
		WithHeader: false,
		bufferSize: int64(os.Getpagesize() * 128),
	}
}

//Split splits file in smaller chunks
func (s Splitter) Split(inputFilePath string, outputDirPath string) ([]string, error) {
	if s.FileChunkSize < minFileChunkSize {
		return nil, ErrSmallFileChunkSize
	}
	file, err := os.Open(inputFilePath)
	if err != nil {
		msg := fmt.Sprintf("Couldn't open file %s : %v", inputFilePath, err)
		return nil, errors.New(msg)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		msg := fmt.Sprintf("Couldn't get file stat %s : %v", inputFilePath, err)
		return nil, errors.New(msg)
	}
	fileSize := stat.Size()
	if fileSize <= s.FileChunkSize {
		return nil, ErrBigFileChunkSize
	}

	bufBulk := make([]byte, s.bufferSize)
	fileName, ext := getFileName(inputFilePath)
	st := ChunkFileInfo{
		s:             s,
		inputFilePath: inputFilePath,
		fileName:      fileName,
		ext:           ext,
		outputDirPath: prepareResultDirPath(outputDirPath),
		inputFile:     file,
		firstLine:     true,
		chunk:         1,
		bulkBuffer:    bytes.NewBuffer(make([]byte, 0, s.bufferSize)),
	}
	for {
		//Read bulk from file
		size, err := file.Read(bufBulk)
		if err == io.EOF {
			_, err = st.chunkFile.Write(st.brokenLine)
			if err != nil {
				msg := fmt.Sprintf("Couldn't write chunk file %s : %v", st.chunkFilePath, err)
				return nil, errors.New(msg)
			}
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Couldn't read file bulk %s : %v", inputFilePath, err)
			return nil, errors.New(msg)
		}
		st.fileBuffer = bytes.NewBuffer(bufBulk[:size])
		if len(st.brokenLine) > 0 {
			st.bulkBuffer.Write(st.brokenLine)
			st.brokenLine = []byte{}
		}

		err = readLinesFromBulk(&st)
		if err != nil {
			return nil, err
		}

		err = saveBulkToFile(&st)
		if err != nil {
			return nil, err
		}
	}
	st.chunkFile.Close()

	return st.result, nil
}

//It reads bulk line by line
func readLinesFromBulk(st *ChunkFileInfo) error {
	for {
		bytesLine, err := st.fileBuffer.ReadBytes('\n')
		if err == io.EOF {
			st.brokenLine = bytesLine
			break
		}
		if err != nil {
			msg := fmt.Sprintf("Couldn't read bytes from buffer of file %s : %v", st.inputFilePath, err)
			return errors.New(msg)
		}
		if st.firstLine && st.s.WithHeader {
			st.firstLine = false
			st.header = bytesLine
			continue
		}
		st.bulkBuffer.Write(bytesLine)
		if st.s.FileChunkSize < st.s.bufferSize && int64(st.bulkBuffer.Len()) >= (st.s.FileChunkSize-int64(len(st.header))) {
			err = saveBulkToFile(st)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// It saves lines from bulk to a new file
func saveBulkToFile(st *ChunkFileInfo) error {
	st.chunkFilePath = st.outputDirPath + st.fileName + "_" + strconv.Itoa(st.chunk) + "." + st.ext
	stat, err := os.Stat(st.chunkFilePath)
	if os.IsNotExist(err) {
		chunkFile, err := os.Create(st.chunkFilePath)
		if err != nil {
			msg := fmt.Sprintf("Couldn't create file %s : %v", st.chunkFilePath, err)
			return errors.New(msg)
		}
		st.setChunkFile(chunkFile)
		//set the file header if it provide the parameter.
		if st.s.WithHeader {
			_, err = st.chunkFile.Write(st.header)
			if err != nil {
				msg := fmt.Sprintf("Couldn't write header of chunk file %s : %v", st.chunkFilePath, err)
				return errors.New(msg)
			}
		}

		st.result = append(st.result, st.chunkFilePath)
	}

	//write the bytes content
	_, err = st.chunkFile.Write(st.bulkBuffer.Bytes())
	if err != nil {
		msg := fmt.Sprintf("Couldn't write chunk file %s : %v", st.chunkFilePath, err)
		return errors.New(msg)
	}
	stat, _ = os.Stat(st.chunkFilePath)
	if stat.Size() > int64(st.s.FileChunkSize-st.s.bufferSize) {
		st.chunk++
	}
	st.bulkBuffer.Reset()
	return nil
}

//It extracts name and extension from path
func getFileName(path string) (string, string) {
	split := strings.Split(path, string(filepath.Separator))
	name := split[len(split)-1]
	split = strings.Split(name, ".")
	ext := split[len(split)-1]
	nameSlice := split[:len(split)-1]
	name = strings.Join(nameSlice, "")

	return name, ext
}

//It handles the output file path.
func prepareResultDirPath(path string) string {
	if path == "" {
		return ""
	}
	p := []byte(path)
	if p[len(p)-1] != filepath.Separator {
		p = append(p, filepath.Separator)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0755); err != nil {
			panic(err)
		}
	}

	return string(p)
}
