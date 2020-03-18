package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	go func() {
		fmt.Println("start pprof")
		http.ListenAndServe("0.0.0.0:9090", nil)
	}()

	s := make([]int, 10000)
	for {
		fill(s)
                sum(s)
	}
}

func fill(s []int) {
	for index := range s {
		s[index] = 1
	}
}

func sum(s []int) int {
	var sum int
	for index := range s {
		sum += s[index]
	}
	return sum
}
