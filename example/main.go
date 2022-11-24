package main

import (
	"fmt"
	"time"
)

var ch = make(chan string)

func main() {

	stratTime := time.Now()

	go count()
	go hesen()

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	fmt.Println(time.Now().Sub(stratTime))
	fmt.Println("Main thread done")
}

func count() {
	time.Sleep(2 * time.Second)
	ch <- "count isledi"
}

func hesen() {
	time.Sleep(2 * time.Second)
	ch <- "hesen isledi"
}
