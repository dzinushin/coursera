package main

import "fmt"

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	select {
	case val1 := <-ch1:
		fmt.Println("ch1 val: ", val1)
	case ch2 <- 2:
		fmt.Println("put val to ch2")
	default:
		fmt.Println("default case")
	}
}
