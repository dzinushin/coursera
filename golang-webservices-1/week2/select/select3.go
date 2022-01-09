package main

import "fmt"

func main() {
	cnclCh := make(chan struct{})
	dataCh := make(chan int)

	go func(cnclCh chan struct{}, dataCh chan int) {
		val := 0
		for {
			select {
			case <-cnclCh:
				return
			case dataCh <- val:
				val++
			}
		}
	}(cnclCh, dataCh)

	for curVal := range dataCh {
		fmt.Println("read: ", curVal)
		if curVal > 3 {
			fmt.Println("send cancel")
			cnclCh <- struct{}{}
			break
		}
	}
}
