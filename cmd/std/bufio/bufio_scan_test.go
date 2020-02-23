package bufio

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"
)

var scanTests = []string{
	"",
	"a",
	"¼",
	"☹",
	"\x81",   // UTF-8 error
	"\uFFFD", // correctly encoded RuneError
	"abcdefgh",
	"abc def\n\t\tgh    ",
	"abc¼☹\x81\uFFFD日本語\x82abc",
}

func TestScanByte(t *testing.T) {
	for n, test := range scanTests {
		buf := strings.NewReader(test)
		s := bufio.NewScanner(buf)
		s.Split(bufio.ScanBytes)
		var i int
		for i = 0; s.Scan(); i++ {
			if b := s.Bytes(); len(b) != 1 || b[0] != test[i] {
				t.Errorf("#%d: %d: expected %q got %q", n, i, test, b)
			}
		}
		if i != len(test) {
			t.Errorf("#%d: termination expected at %d; got %d", n, len(test), i)
		}
		err := s.Err()
		if err != nil {
			t.Errorf("#%d: %v", n, err)
		}
	}
}

func TestBufIO(t *testing.T) {
	filename := "temp.txt"
	file, err := os.Create(filename)
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	writer := bufio.NewWriterSize(file, 20)
	linesToWrite := []string{"Rust--A language empowering everyone to build reliable and efficient software. ",
		"Rust’s rich type system and ownership ",
		"Rust has great documentation, a friendly compiler with useful error messages, and top-notch tooling — an integrated",
		"Rust is blazingly fast and memory-efficient"}
	for _, line := range linesToWrite {
		bytesWritten, err := writer.WriteString(line + "\n")
		if err != nil {
			t.Fatalf("Got error while writing to a file. Err: %s", err.Error())
		}
		fmt.Printf("Bytes Written: %d\n", bytesWritten)
		fmt.Printf("Available: %d\n", writer.Available())
		fmt.Printf("Buffered : %d\n", writer.Buffered())
	}
	writer.Flush()
	os.Remove(filename)

}
