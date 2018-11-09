package concurrency

import (
	"fmt"
	"log"
	"testing"
)

func TestParallelize(t *testing.T) {
	//prepare data
	data := make([][]int, 10)
	for k, _ := range data {
		data[k] = []int{k + 10, k + 20}
	}
	log.Printf("\n%+v\n", data)

	fss := make([]func() interface{}, 0)
	for _, v := range data {

		r := func(a, b int) func() interface{} {
			return func() interface{} {
				return add(a, b)
			}
		}(v[0], v[1])
		fss = append(fss, r)
	}

	result := Parallelize(fss)
	fmt.Printf("\n%+v\n", result)
	//test if we got expected output.
	expected := []interface{}{
		30, 32, 34, nil, 38, 40, 42, 44, 46, 48,
	}
	for k, _ := range result {
		if result[k] != expected[k] {
			t.Errorf("Expected %v, Got %v at Index: %d", expected[k], result[k], k)
		}
	}
}

func TestParallelizeThrottled(t *testing.T) {
	//prepare data
	data := make([][]int, 10)
	for k, _ := range data {
		data[k] = []int{k + 10, k + 20}
	}
	log.Printf("\n%+v\n", data)

	fss := make([]func() interface{}, 0)
	for _, v := range data {

		r := func(a, b int) func() interface{} {
			return func() interface{} {
				return add(a, b)
			}
		}(v[0], v[1])
		fss = append(fss, r)
	}

	result := ParallelizeThrottled(fss, 3)
	fmt.Printf("\n%+v\n", result)
	//test if we got expected output.
	expected := []interface{}{
		30, 32, 34, nil, 38, 40, 42, 44, 46, 48,
	}
	for k, _ := range result {
		if result[k] != expected[k] {
			t.Errorf("Expected %v, Got %v at Index: %d", expected[k], result[k], k)
		}
	}
}

func add(a, b int) int {
	if a == 13 {
		panic("p[an")
	}
	return a + b
}
