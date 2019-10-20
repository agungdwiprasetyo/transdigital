package inventory

import (
	"agungdwiprasetyo/transdigital/task2/shared"
	"fmt"
	"log"
	"sync"
	"time"
)

// Inventory model
type Inventory struct {
	handlers []shared.Handler
	prefix   string
}

// NewInventory constructor
func NewInventory() *Inventory {
	prefix := fmt.Sprintf("\x1b[32;1mInventory\x1b[0m")
	return &Inventory{
		handlers: []shared.Handler{
			newHandler(prefix, "REST Handler"), newHandler(prefix, "RPC Handler"),
		},
		prefix: fmt.Sprintf("\x1b[32;1mInventory\x1b[0m"),
	}
}

// Run Inventory worker
func (c *Inventory) Run(wg *sync.WaitGroup, job chan string) {
	wg.Add(1)
	defer wg.Done()

	for j := range job {

		log.Printf("%s: Worker started job %s\n", c.prefix, j)

		messageFromWorker := fmt.Sprintf("%s:%s", j, shared.GenerateMessage(15))
		time.Sleep(1 * time.Second) // latency from worker

		for _, handler := range c.handlers {
			var subWg sync.WaitGroup
			handlerJob := make(chan shared.MessageFormat, shared.MaxThreads)
			go handler.Handler(&subWg, handlerJob)

			// Workers should send jobs over channels to each subordinate. Each subordinate will have 20 jobs.
			for i := 1; i <= shared.MaxSubordinateJob; i++ {
				handlerJob <- shared.MessageFormat{JobID: i, Message: messageFromWorker}
			}
			close(handlerJob)
			subWg.Wait()
		}

		log.Printf("%s: Worker finished job %s with result %s\n", c.prefix, j, messageFromWorker)
	}
}

type handler struct {
	workerName string
	protocol   string
}

func newHandler(workerName, protocol string) *handler {
	return &handler{workerName, protocol}
}

func (h *handler) Handler(wg *sync.WaitGroup, messages chan shared.MessageFormat) {
	wg.Add(1)
	defer wg.Done()

	for message := range messages {
		time.Sleep(100 * time.Millisecond) // latency from handler send message
		log.Printf("%s: %s finish send message [%s] from worker with job ID %d\n", h.workerName, h.protocol, message.Message, message.JobID)
	}
}
