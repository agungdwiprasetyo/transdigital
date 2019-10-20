package main

// Handling race condition (case checkout same product in same time) using mutual exclusion

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

	// mutual exclusion for handle race condition if checkout product in same time
	sync.Mutex
}

type order struct {
	SKU   string
	Total int
}

type user struct {
	ID    string
	Order *order
}

// global stocks in store database (list of product object)
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

func checkout(ord *order) error {
	p := findSkuInProduct(ord.SKU)
	if p == nil {
		return errors.New("product not found")
	}

	// lock for handle if this function call with goroutine
	p.Lock()
	defer p.Unlock()

	if p.Stock == 0 {
		return fmt.Errorf("product %s out of stock", ord.SKU)
	}

	time.Sleep(3 * time.Second) // assume heavy process when success take product

	// validate stock and total order
	if p.Stock < ord.Total {
		return errors.New("insufficient stock")
	}

	// update product stock
	p.Stock -= ord.Total

	return nil
}

func handler(u *user) {
	fmt.Printf("%s: handling request from %s to checkout product %s\n", time.Now().Format(time.RFC3339Nano), u.ID, u.Order.SKU)
	err := checkout(u.Order)
	if err != nil {
		fmt.Printf("\x1b[31;1m%s can't checkout SKU %s, Error: %v\x1b[0m\n", u.ID, u.Order.SKU, err)
	} else {
		fmt.Printf("\x1b[32;1mProduct %s success taken by %s\x1b[0m\n", u.Order.SKU, u.ID)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill)

	sku := "SKU001"
	var activeUser []*user

	userA := &user{
		ID: "UserA",
		Order: &order{
			SKU:   sku,
			Total: 1,
		},
	}
	activeUser = append(activeUser, userA)
	userB := &user{
		ID: "UserB",
		Order: &order{
			SKU:   sku,
			Total: 1,
		},
	}
	activeUser = append(activeUser, userB)
	userC := &user{
		ID: "UserC",
		Order: &order{
			SKU:   sku,
			Total: 1,
		},
	}
	activeUser = append(activeUser, userC)

	// user A and user B (and additional user C) checkout in same time
	for _, u := range activeUser {
		go handler(u)
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
