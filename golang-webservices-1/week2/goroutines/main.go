package main

import (
	"fmt"
	"runtime"
)

func main() {
	for i := 0; i < 4; i++ {
		go doSomeWork(i)
	}
	fmt.Scanln()
}

func doSomeWork(n int) {
	for i := 0; i < 10; i++ {
		fmt.Printf("th: %v work: %v\n", n, i)
		runtime.Gosched()
	}
}
