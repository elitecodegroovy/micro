package main

import (
	"flag"
	"github.com/micro/micro/cmd/file/bufv2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

const (
	tinyFileRegexPattern string = `\w+_\d\.\w+$`
)

var (
	regexFileName = regexp.MustCompile(tinyFileRegexPattern)
)

func main() {
	var files = flag.String("files", "", "provide the file chunk size for splitter (B unit), default 10485760")
	flag.Parse()

	var filenames []string
	if len(*files) == 0 {
		files, err := ioutil.ReadDir(".")
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if regexFileName.MatchString(f.Name()) {
				filenames = append(filenames, f.Name())
				log.Printf("filename: %s", f.Name())
			}

		}
	} else {
		filenames = strings.Split(*files, ",")
	}

	for i := 0; i < len(filenames); i++ {
		log.Printf("filename: %s", filenames[i])
		buf.ShowMemoryInfo()
		buf.ReadLinesByBufIO(filenames[i])
		buf.Done()
	}

}
