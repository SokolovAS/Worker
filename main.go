package main

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

type Worker struct {
	min, max int
}

func worker(ctx context.Context, count *int32, tasksCh <-chan func(), w Worker) {
	atomic.AddInt32(count, 1)

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
			fmt.Println("int32(min)", int32(w.min))
			if atomic.LoadInt32(count) <= int32(w.min) {
				continue
			} else {
				atomic.AddInt32(count, -1)
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func NewWP(ctx context.Context, ch chan func()) {
	var actualWorkersCount int32
	w := Worker{3, 10}

	tasksCh := make(chan func())
	for {
		select {
		case t := <-ch:
			select {
			case tasksCh <- t:
			case <-time.After(100 * time.Millisecond):
				fmt.Println("actualWorkersCount", atomic.LoadInt32(&actualWorkersCount))
				if atomic.LoadInt32(&actualWorkersCount) < int32(w.max) {
					fmt.Println("make worker")
					go worker(ctx, &actualWorkersCount, tasksCh, w)
				}
				tasksCh <- t
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	tasks := 100
	ch := make(chan func())
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	go NewWP(ctx, ch)
	for i := 0; i < tasks; i++ {
		ch <- func() {
			time.Sleep(time.Second)
		}
	}
	time.Sleep(3 * time.Second)
	cancel()
}
