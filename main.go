package main

import (
	"fmt"
	"sync"
	"time"
)

const defaultWorkersCount = 2

func worker(tasksCh <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		task, ok := <-tasksCh
		if !ok {
			return
		}
		d := time.Duration(task) * time.Millisecond
		time.Sleep(d)
		fmt.Println("processing task", task)
	}
}

func pool(wg *sync.WaitGroup, dch chan int, lch chan int) {
	tasksCh := make(chan int)

	for {
		select {
		case task := <-lch:
			fmt.Println("in additional load")
			go worker(tasksCh, wg)
			tasksCh <- task
		case task := <-dch:
			fmt.Println("in default loads")
			go worker(tasksCh, wg)

			tasksCh <- task

			wg.Done()
		}
	}
}

func increaseLoads(tasks int, ch chan int) {
	for i := 37; i < tasks; i++ {
		ch <- i
	}
}

func defaultLoads(tasks int, ch chan int) {
	for i := 0; i < tasks; i++ {
		ch <- i
	}
}

func main() {
	var defaultTasks = 36
	var additionalTasks = 57

	highLoadCh := make(chan int)
	defaultLoadCh := make(chan int)

	t := time.Now()
	var wg sync.WaitGroup
	wg.Add(defaultWorkersCount + additionalTasks)
	go pool(&wg, defaultLoadCh, highLoadCh)
	go defaultLoads(defaultTasks, defaultLoadCh)
	go increaseLoads(additionalTasks, highLoadCh)
	wg.Wait()
	t1 := time.Now()
	dif := t1.Sub(t)
	fmt.Println(dif)
}
