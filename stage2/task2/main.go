package main

// limit goroutine using buffered channel

import (
	"agungdwiprasetyo/transdigital/task2/customer"
	"agungdwiprasetyo/transdigital/task2/inventory"
	"agungdwiprasetyo/transdigital/task2/logistic"
	"agungdwiprasetyo/transdigital/task2/order"
	"agungdwiprasetyo/transdigital/task2/product"
	"agungdwiprasetyo/transdigital/task2/shared"
	"fmt"
	"strconv"
	"sync"
	"time"
)

type worker struct {
	customer shared.Worker
}

func main() {

	workers := []shared.Worker{
		customer.NewCustomer(), inventory.NewInventory(), logistic.NewLogistic(), order.NewOrder(), product.NewProduct(),
	}

	for i, worker := range workers {
		now := time.Now()
		var wg sync.WaitGroup
		jobs := make(chan string, shared.MaxThreads)

		for i := 0; i < shared.MaxThreads; i++ {
			go worker.Run(&wg, jobs)
		}

		// Main thread should send jobs over channels to each worker. Each worker will have 10 jobs.
		for j := 1; j <= shared.MaxWorkerJob; j++ {
			jobs <- strconv.Itoa(j)
		}
		close(jobs)

		// wait current worker finish processing jobs
		wg.Wait()
		fmt.Printf("\x1b[31;1mWORKER %d FINISHED IN %v\x1b[0m\n", i+1, time.Since(now))
	}
}
