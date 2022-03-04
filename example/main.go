package main

import (
	"fmt"
	"forrest"
	"sync"
	"sync/atomic"
	"time"
)

var sum int32

func myFunc(i interface{}) {
	n := i.(int32)
	atomic.AddInt32(&sum, n)
	fmt.Printf("run with %d\n", n)
}

func demoFunc() {
	time.Sleep(10 * time.Millisecond)
	fmt.Println("Hello World!")
}

func main() {
	runTimes := 1000

	var wg sync.WaitGroup
	syncCalculateSum := func() {
		demoFunc()
		wg.Done()
	}
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		_ = forrest.Submit(syncCalculateSum)
	}
	wg.Wait()
	fmt.Printf("running goroutines: %d\n", forrest.Running())
	fmt.Printf("finish all tasks.\n")

	p2 := forrest.NewPool(2, 10)
	for i := 0; i < runTimes; i++ {
		wg.Add(1)
		j := i
		_ = p2.Submit(func() {
			myFunc(int32(j))
			wg.Done()
		})
	}

	wg.Wait()
	fmt.Printf("running goroutines: %d\n", p2.Running())
	fmt.Printf("finish all tasks, result is %d\n", sum)
	if sum != 499500 {
		panic("the final result is wrong!!!")
	}
}
