package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func generateSlice(size int) []int64 {

	slice := make([]int64, size, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Int63() - rand.Int63()
	}
	return slice
}

func TestSort(t *testing.T) {
	a1 := generateSlice(30)
	fmt.Println("\n--- Unsorted --- ", a1)
	selectionSort(a1)
	fmt.Println("\n--- Unsorted --- ", a1)
	for i, v := range a1 {
		fmt.Printf("i=%d, v=%d \n", i, v)
	}
}
