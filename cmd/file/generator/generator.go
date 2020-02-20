package main

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	fileName = "bigLongTypeData.txt"
)

const newl = "\n"

func init() {
	//以时间作为初始化种子
	rand.Seed(time.Now().UnixNano())
}

//Automatically generate int64 and int32 number data.
func generateRandomLongNum() {
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		fmt.Printf("failed to open file : %s ", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	var g errgroup.Group

	// 0.900 G data will be generated.
	for i := 0; i < 1024; i++ {
		for j := 0; j < 1024; j++ {
			g.Go(func() error {
				var v string
				if i%2 == 0 {
					v = strconv.FormatInt(int64(rand.Int31()), 10)
				} else {
					v = strconv.FormatInt(rand.Int63(), 10)
				}

				fmt.Printf("i(%d) j(%d) value :%s \n", i, j, v)
				f.Write([]byte(v + newl))
				return nil
			})
		}
		// Wait for all HTTP fetches to complete.
		if err := g.Wait(); err == nil {
			fmt.Printf(">>>>>>>>>>>>>>>>>>>> i: %d \n", i)
			if fileInfo, err := f.Stat(); err == nil {
				//900M
				fmt.Printf(">>>>>>>>>>>>>>>>>>>>　file size:: %f M \n", float64(fileInfo.Size())/(float64(1024)*float64(1024)))
			}
		}
	}
}

func main() {
	//step 1
	generateRandomLongNum()

}