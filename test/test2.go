package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(id int, exitCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		fmt.Printf("\nIn worker for loop\n")
		select {
		case <-exitCh:
			fmt.Printf("Worker %d received exit signal. Exiting.\n", id)
			return
		default:
			// Your worker logic here
			fmt.Printf("Worker %d is working...\n", id)
		}
		time.Sleep(2*time.Second)
	}
}

func main() {
	exitCh := make(chan struct{})
	var wg sync.WaitGroup

	// Start three worker goroutines
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go worker(i, exitCh, &wg)
	}

	// Simulate some work before sending the exit signal
	// In a real scenario, you might perform other tasks here

	// Send the exit signal to all workers
	fmt.Printf("\nclosing exitCh now\n")
	//close(exitCh)
	exitCh <- struct{}{}

	// Wait for all worker goroutines to finish
	wg.Wait()

	fmt.Println("All workers have exited.")
}
