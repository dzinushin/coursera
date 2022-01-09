package main

import "fmt"

func main() {
	ch := make(chan int)

	go func(in chan int) {
		fmt.Println("GOROUTINE: before read from chan")
		val := <-in
		fmt.Println("GOROUTINE: get value from chan: ", val)
	}(ch)

	fmt.Println("MAIN: before put to chan")
	ch <- 42
	ch <- 44

	fmt.Println("MAIN: after put to chan")
	fmt.Scanln()
}
