package main

import "fmt"

func main() {
	ch1 := make(chan int, 2)
	ch1 <- 10
	ch1 <- 11
	ch2 := make(chan int, 2)
	ch2 <- 20

	fmt.Println("before for")
LOOP:
	for {
		fmt.Println("in for")
		select {
		case v1 := <-ch1:
			fmt.Println("ch1 val: ", v1)
		case v2 := <-ch2:
			fmt.Println("ch2 val: ", v2)
		default:
			fmt.Println("default")
			break LOOP
		}
	}
	fmt.Println("after for")
}
