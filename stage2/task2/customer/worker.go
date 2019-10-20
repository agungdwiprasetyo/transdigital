package customer

import (
	"agungdwiprasetyo/transdigital/task2/shared"
	"fmt"
	"log"
	"sync"
	"time"
)

type messageFormat struct {
	JobID   int
	Message string
}

// Customer model
type Customer struct {
	prefix string
}

// NewCustomer constructor
func NewCustomer() *Customer {
	return &Customer{
		prefix: fmt.Sprintf("\x1b[32;1mCustomer\x1b[0m"),
	}
}

// Run customer worker
func (c *Customer) Run(wg *sync.WaitGroup, job chan string) {
	wg.Add(1)
	defer wg.Done()

	for j := range job {

		log.Printf("%s: Worker started job %s\n", c.prefix, j)

		messageFromWorker := fmt.Sprintf("%s:%s", j, shared.GenerateMessage(15))
		time.Sleep(1 * time.Second) // latency from worker

		handlerJob := make(chan messageFormat, shared.MaxThreads)
		go c.RestHandler(handlerJob)
		go c.RPCHandler(handlerJob)

		// Workers should send jobs over channels to each subordinate. Each subordinate will have 20 jobs.
		for i := 1; i <= shared.MaxSubordinateJob; i++ {
			handlerJob <- messageFormat{JobID: i, Message: messageFromWorker}
		}
		close(handlerJob)

		log.Printf("result from worker %s with job %s is %s\n", c.prefix, j, messageFromWorker)
	}
}

// RestHandler customer
func (c *Customer) RestHandler(messages chan messageFormat) {
	for message := range messages {
		time.Sleep(500 * time.Millisecond) // latency from rest handler
		log.Printf("%s: RestHandler finish send message %s from worker with job ID %d\n", c.prefix, message.Message, message.JobID)
	}
}

// RPCHandler customer
func (c *Customer) RPCHandler(messages chan messageFormat) {
	for message := range messages {
		time.Sleep(500 * time.Millisecond) // latency from rpc handler
		log.Printf("%s: RPCHandler finish send message %s from worker with job ID %d\n", c.prefix, message.Message, message.JobID)
	}
}
