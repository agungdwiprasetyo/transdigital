package shared

import "sync"

type Worker interface {
	Run(wg *sync.WaitGroup, job chan string)
}
