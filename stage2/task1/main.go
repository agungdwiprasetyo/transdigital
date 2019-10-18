package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type product struct {
	SKU   string
	Stock int
}

// global stocks in store (list of product object)
var products = []*product{
	{
		SKU:   "SKU001",
		Stock: 1,
	},
	{
		SKU:   "SKU002",
		Stock: 2,
	},
}

func findSkuInProduct(sku string) *product {
	for _, p := range products {
		if p.SKU == sku {
			return p
		}
	}
	return nil
}

// lock global mutual exclusion
var lock sync.Mutex

func checkout(skuNo string) error {
	// lock for handle race condition if checkout in same time
	lock.Lock()
	defer lock.Unlock()

	p := findSkuInProduct(skuNo)
	if p == nil {
		return errors.New("product not found")
	}

	if p.Stock == 0 {
		return fmt.Errorf("product %s out of stock", skuNo)
	}

	time.Sleep(3 * time.Second) // assume heavy process when success take product

	// update product stock
	p.Stock--

	return nil
}

func handler(sku string, result chan error) {
	fmt.Printf("%s: handling request checkout sku %s\n", time.Now().Format(time.RFC3339Nano), sku)
	result <- checkout(sku)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	userA := make(chan error)
	userB := make(chan error)
	userC := make(chan error)
	sku := "SKU001"

	// user A and user B (and user C) checkout in same time
	go handler(sku, userA)
	go handler(sku, userB)
	go handler(sku, userC)

	// take result each channel (user)
	if err := <-userA; err != nil {
		fmt.Printf("User A can't checkout SKU %s, %v\n", sku, err)
	}
	if err := <-userB; err != nil {
		fmt.Printf("User B can't checkout SKU %s, %v\n", sku, err)
	}
	if err := <-userC; err != nil {
		fmt.Printf("User C can't checkout SKU %s, %v\n", sku, err)
	}

	// wait timeout or interupt or kill signal channel
	select {
	case <-quit:
		log.Println("os interrupted")
	case <-ctx.Done():
		log.Println("context timeout")
	}

	p := findSkuInProduct(sku)
	if p.Stock < 0 {
		panic(fmt.Sprintf("\x1b[31;1mWRONG ANSWER SOLUTION, STOCK IS LESS THAN ZERO (%d)\x1b[0m", p.Stock))
	}
	fmt.Printf("\x1b[32;1mOK\x1b[0m\n")
}
