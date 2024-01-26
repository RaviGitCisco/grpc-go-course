package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int)

	var wg sync.WaitGroup
	wg.Add(3)

	// Three goroutines receiving from the same channel
	go func() {
		defer wg.Done()
		val := <-ch
		fmt.Println("Goroutine 1 received:", val)
	}()

	go func() {
		defer wg.Done()
		val := <-ch
		fmt.Println("Goroutine 2 received:", val)
	}()

	go func() {
		defer wg.Done()
		val := <-ch
		fmt.Println("Goroutine 3 received:", val)
	}()

	// Sending a value to the channel
	ch <- 42

	// Wait for all goroutines to finish
	wg.Wait()
}
