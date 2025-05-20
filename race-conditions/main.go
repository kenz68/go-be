package main

import (
	"fmt"
	"sync"
)

var counter int

func increment(wg *sync.WaitGroup) {
	counter++
	wg.Done()
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go increment(&wg)
	go increment(&wg)

	wg.Wait()
	fmt.Println("Final counter value:", counter)
}
