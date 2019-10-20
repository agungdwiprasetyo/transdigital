package main

// limit goroutine using buffered channel

import (
	"agungdwiprasetyo/transdigital/task2/customer"
	"agungdwiprasetyo/transdigital/task2/shared"
	"strconv"
	"sync"
)

type worker struct {
	customer shared.Worker
}

func main() {
	jobs := make(chan string, shared.MaxThreads)

	var w = worker{
		customer: customer.NewCustomer(),
	}

	var wg sync.WaitGroup
	for i := 0; i < shared.MaxThreads; i++ {
		go w.customer.Run(&wg, jobs)
	}

	// Main thread should send jobs over channels to each worker. Each worker will have 10 jobs.
	for j := 1; j <= shared.MaxWorkerJob; j++ {
		jobs <- strconv.Itoa(j)
	}
	close(jobs)

	// wait each worker finish process jobs
	wg.Wait()
}
