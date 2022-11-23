package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	stratTime := time.Now()

	var wg sync.WaitGroup

	wg.Add(2)
	go count(1, &wg)
	go hesen(2, &wg)

	wg.Wait()

	fmt.Println(time.Now().Sub(stratTime))
	fmt.Println("Main thread done")
}

func count(a int, wg *sync.WaitGroup) {

	/*for i := 0; i < 10; i++ {
		fmt.Println("Thread a=", a, " ", i)
	}*/

	time.Sleep(3 * time.Second)
	wg.Done()

}

func hesen(a int, wg *sync.WaitGroup) {

	/*for i := 0; i < 10; i++ {
		fmt.Println("Thread a=", a, " ", i)
	}*/

	time.Sleep(4 * time.Second)
	wg.Done()

}
