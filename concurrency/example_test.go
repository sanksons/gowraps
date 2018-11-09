package concurrency

import (
	"fmt"
)

func ExampleParallelize_sum() {

	// Example:
	// We need to calculate sum of two numbers
	// but we have 10 sets of such numbers
	//
	data := [][]int{
		[]int{10, 20}, []int{11, 21}, []int{12, 22}, []int{13, 23}, []int{14, 24},
		[]int{15, 25}, []int{16, 26}, []int{17, 27}, []int{18, 28}, []int{19, 29},
	}
	//Following is the sum function that we will call in parallel.
	sum := func(a, b int) int {
		return a + b
	}

	//Since, our parallelize func accepts []func() interface{}
	// we need to transform our dataset in that format before passing.
	var fss []func() interface{}
	for _, v := range data {
		r := func(a, b int) func() interface{} {
			return func() interface{} {
				return sum(a, b)
			}
		}(v[0], v[1])
		fss = append(fss, r)
	}

	result := Parallelize(fss)
	fmt.Printf("%+v", result)
	// Output: [30 32 34 36 38 40 42 44 46 48]
}

func ExampleParallelizeThrottled_sum() {

	// Example:
	// We need to calculate sum of two numbers
	// but we have 10 sets of such numbers and we will limit parallelism to 3.
	//
	data := [][]int{
		[]int{10, 20}, []int{11, 21}, []int{12, 22}, []int{13, 23}, []int{14, 24},
		[]int{15, 25}, []int{16, 26}, []int{17, 27}, []int{18, 28}, []int{19, 29},
	}
	//Following is the sum function that we will call in parallel.
	sum := func(a, b int) int {
		return a + b
	}

	//Since, our parallelize func accepts []func() interface{}
	// we need to transform our dataset in that format before passing.
	var fss []func() interface{}
	for _, v := range data {
		r := func(a, b int) func() interface{} {
			return func() interface{} {
				return sum(a, b)
			}
		}(v[0], v[1])
		fss = append(fss, r)
	}

	result := ParallelizeThrottled(fss, 3)
	fmt.Printf("%+v", result)
	// Output: [30 32 34 36 38 40 42 44 46 48]
}
