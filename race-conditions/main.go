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

var mutex sync.Mutex

// This function uses a mutex to safely increment the counter
func incrementMutex(wg *sync.WaitGroup) {
	mutex.Lock()
	counter++
	mutex.Unlock()
	wg.Done()
}
func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go incrementMutex(&wg)
	go incrementMutex(&wg)

	wg.Wait()
	fmt.Println("Final counter value:", counter)
}
