package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func worker(count *int32, tasksCh <-chan func()) {
	min := 2

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		atomic.AddInt32(count, 1)
		wg.Done()
	}()
	wg.Wait()

	tick := time.Tick(2 * time.Second)

	for {
		select {
		case t := <-tasksCh:
			fmt.Println("processing task")
			t()
			fmt.Println("Done!")
		case <-tick:
			fmt.Println("Tick!")
			fmt.Println("atomic.LoadInt32(count)", atomic.LoadInt32(count))
			fmt.Println("int32(min)", int32(min))
			if atomic.LoadInt32(count) <= int32(min) {
				continue
			} else {
				//need decrease count
				return
			}
		}
	}
}

func NewWP(ch chan int) {
	var actualWorkersCount int32
	max := 10

	tasksCh := make(chan func())

	for {
		select {
		case <-ch:
			fmt.Println("actualWorkersCount", atomic.LoadInt32(&actualWorkersCount))
			if atomic.LoadInt32(&actualWorkersCount) < int32(max) {
				fmt.Println("make worker")
				go worker(&actualWorkersCount, tasksCh)
				tasksCh <- func() {
					time.Sleep(time.Second)
				}
			}
			continue
		}
	}
}

func main() {
	tasks := 100

	ch := make(chan int)
	go NewWP(ch)

	for i := 0; i < tasks; i++ {
		//time.Sleep(100 * time.Millisecond)
		ch <- i
	}

	time.Sleep(30 * time.Second)
}
