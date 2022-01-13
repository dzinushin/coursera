package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func startJob(in, out chan interface{}) {
	fmt.Println("startJob start")
	//data := []int{0, 1, 1, 2, 3, 5, 8}
	data := []int{0, 1}
	for i := range data {
		out <- i
	}
	close(out)
}

func finalJob(in, out chan interface{}) {
	fmt.Println("finalJob start")
	cnt := 0
	for data := range in {
		cnt++
		fmt.Println("finalJob received data", data)
	}
	fmt.Println("finalJob received ", cnt, "items")
	close(out)
}

func main() {
	jobs := []job{
		startJob,
		SingleHash,
		MultiHash,
		CombineResults,
		finalJob,
	}
	ExecutePipeline(jobs...)
}

func ExecutePipeline(jobs ...job) {
	fmt.Println(len(jobs))
	in := make(chan interface{})

	wg := &sync.WaitGroup{}
	for _, theJob := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go jobRunner(theJob, in, out, wg)
		in = out
	}

	wg.Wait()
}

func jobRunner(theJob job, in, out chan interface{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()
	theJob(in, out)
	close(out)
}

func SingleHash(in, out chan interface{}) {
	fmt.Println("SingleHash start")
	wg := &sync.WaitGroup{}
	m := &sync.Mutex{}
	for data := range in {
		fmt.Println("SingleHash received data", data)
		s := strconv.Itoa(data.(int))
		wg.Add(1)
		go singleHashCalc(s, m, wg, out)
	}

	wg.Wait()
}

func singleHashCalc(data string, m *sync.Mutex, wg *sync.WaitGroup, out chan interface{}) {
	defer func() { wg.Done() }()

	m.Lock()
	md5sum := DataSignerMd5(data)
	m.Unlock()

	ch1 := make(chan string)
	go crc32Calc(md5sum, ch1)

	ch2 := make(chan string)
	go crc32Calc(data, ch2)

	crc32md5sum := <-ch1
	crc32sum := <-ch2

	result := crc32sum + "~" + crc32md5sum
	fmt.Println("SingleHash result", result)
	out <- result
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for data := range in {
		dataStr := data.(string)
		wg.Add(1)
		go multiHashCalc(dataStr, out, wg)
	}

	wg.Wait()
}

func multiHashCalc(data string, out chan interface{}, wg *sync.WaitGroup) {
	defer func() { wg.Done() }()
	var channels []chan string
	for th := 0; th < 6; th++ {
		ch := make(chan string)
		channels = append(channels, ch)
		go crc32Calc(strconv.Itoa(th)+data, ch)
	}

	var result string
	for _, ch := range channels {
		crc32 := <-ch
		result += crc32
	}

	fmt.Printf("%s MutliHash: result: %s\n", data, result)
	out <- result
}

func crc32Calc(data string, out chan string) {
	crc32 := DataSignerCrc32(data)
	out <- crc32
	close(out)
}

func CombineResults(in, out chan interface{}) {
	inputs := make([]string, 0)
	for data := range in {
		fmt.Println("CombineResults: received data ", data)
		inputs = append(inputs, data.(string))
	}

	sort.Strings(inputs)

	result := strings.Join(inputs, "_")
	fmt.Println("CombineResults: result ", result)
	out <- result
}
