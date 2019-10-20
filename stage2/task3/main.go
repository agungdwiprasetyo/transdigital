package main

import (
	"fmt"
	"net"
	"sync"
)

func handler(c net.Conn, ch chan string) {
	addr := c.RemoteAddr().Network()
	ch <- addr
	c.Write([]byte(`{"response":"OK"}`))
	c.Close()
}
func send(sch chan int, results chan int) {
	for {
		data := <-sch
		data++
		results <- data
	}
}
func enabler(results chan int) {
	for {
		<-results
	}
}
func pub(ch chan string, n int) {
	sch := make(chan int)
	results := make(chan int)
	for i := 0; i < n; i++ {
		go send(sch, results)
	}
	go enabler(results)
	for {
		addr := <-ch
		l := len(addr)
		sch <- l
	}
}
func sub(l net.Listener, ch chan string) {
	for {
		c, err := l.Accept()
		if err != nil {
			continue
		}
		go handler(c, ch)
	}
}
func main() {
	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		panic(err)
	}
	ch := make(chan string)

	// menggunakan wait group untuk menghindari os exit ketika menjalankan beberapa goroutine pada fungsi main
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		pub(ch, 5)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sub(l, ch)
	}()

	// wait all goroutine
	wg.Wait()
	fmt.Println("exit")
}
