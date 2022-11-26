package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	stratTime := time.Now()

	var ch1 = make(chan string)
	var ch2 = make(chan string)

	go func() {
		ch1 <- "hi"
	}()

	go func() {
		ch2 <- "there"
	}()

	for {
		select {
		case msg := <-ch1:
			fmt.Println(msg)
		case msg := <-ch2:
			fmt.Println(msg)
			os.Exit(0)
		}
	}

	fmt.Println(time.Now().Sub(stratTime))
	fmt.Println("Main thread done")
}

/*func count() {
	time.Sleep(2 * time.Second)
	ch <- "count isledi"
}

func hesen() {
	time.Sleep(2 * time.Second)
	ch <- "hesen isledi"
}
*/
