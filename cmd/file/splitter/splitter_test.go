package splitter

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var filesDefaultFlow = []string{
	"data/result_default/bigLongTypeData_1.txt",
	"data/result_default/bigLongTypeData_2.txt",
	"data/result_default/bigLongTypeData_3.txt",
	"data/result_default/bigLongTypeData_4.txt",
	"data/result_default/bigLongTypeData_5.txt",
}
var filesWithoutHeader = []string{
	"data/result_without_header/bigLongTypeData_1.txt",
	"data/result_without_header/bigLongTypeData_2.txt",
	"data/result_without_header/bigLongTypeData_3.txt",
	"data/result_without_header/bigLongTypeData_4.txt",
	"data/result_without_header/bigLongTypeData_5.txt",
}
var filesSmallBuffer = []string{
	"data/result_small_buffer/bigLongTypeData_1.txt",
	"data/result_small_buffer/bigLongTypeData_2.txt",
	"data/result_small_buffer/bigLongTypeData_3.txt",
	"data/result_small_buffer/bigLongTypeData_4.txt",
	"data/result_small_buffer/bigLongTypeData_5.txt",
}

func setUp(t *testing.T) {
	files := append(filesDefaultFlow, filesWithoutHeader...)
	files = append(files, filesSmallBuffer...)
	for _, file := range files {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			continue
		}
		err = os.Remove(file)
		if err != nil {
			t.Error(err)
		}
	}
}

func TestSplitFile(t *testing.T) {
	setUp(t)
	input := "data/bigLongTypeData.txt"
	t.Run("Default flow", func(t *testing.T) {
		t1 := time.Now()
		s := New()
		s.FileChunkSize = 10485760
		result, err := s.Split(input, "data/result_default")
		assertResult(t, result, filesDefaultFlow)
		assert.Nil(t, err)
		t.Logf("time elapsed %f s", time.Since(t1).Seconds())
	})
	t.Run("Without headers", func(t *testing.T) {
		s := New()
		s.FileChunkSize = 10485760
		s.WithHeader = false
		result, err := s.Split(input, "data/result_without_header")
		assertResult(t, result, filesWithoutHeader)
		assert.Nil(t, err)
	})
	t.Run("With small buffer", func(t *testing.T) {
		s := New()
		s.FileChunkSize = 10485760
		s.bufferSize = 100
		_, err := s.Split(input, "data/result_small_buffer/")
		//assertResult(t, result, filesSmallBuffer)
		assert.Nil(t, err)
	})
	t.Run("Big file chunk", func(t *testing.T) {
		s := New()
		s.FileChunkSize = 10485760000000
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.Equal(t, err, ErrBigFileChunkSize)
	})
	t.Run("Small file chunk error", func(t *testing.T) {
		s := New()
		result, err := s.Split(input, "")

		assert.Nil(t, result)
		assert.Equal(t, err, ErrSmallFileChunkSize)
	})
	setUp(t)
}

func assertResult(t *testing.T, result []string, expected []string) {
	assert.Equal(t, expected, result)
}

func TestSingleQuota(t *testing.T) {
	fmt.Println(">>>>[", string(filepath.Separator), "]")
}
