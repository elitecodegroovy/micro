package main

import (
	"flag"
	"github.com/micro/micro/cmd/file/mergence"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

const (
	tinyFileRegexPattern string = `\w+_\d\w+\.\w+$`
)

var (
	regexFileName = regexp.MustCompile(tinyFileRegexPattern)
)

var sortedFilename = flag.String("sorted-filename", "data/sortedBigLongTypeData.txt", "provide the output filename, default sorted data file filename 'data/sortedBigLongTypeData.txt'")

func main() {
	flag.Parse()

	var filenames []string
	files, err := ioutil.ReadDir("./data/")
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if regexFileName.MatchString(f.Name()) {
			filename := "data" + string(filepath.Separator) + f.Name()
			filenames = append(filenames, filename)
			log.Printf("filename: %s", f.Name())
		}
	}
	//
	if len(filenames) != 0 {
		m := mergence.New(*sortedFilename, 1048576, filenames)
		if err = m.Merge(); err != nil {
			log.Fatal("..." + err.Error())
		}
	}

}
