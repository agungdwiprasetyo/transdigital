package shared

import "sync"

type Worker interface {
	Run(wg *sync.WaitGroup, job chan string)
}

type Handler interface {
	Handler(wg *sync.WaitGroup, messages chan MessageFormat)
}
