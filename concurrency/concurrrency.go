package concurrency

import (
	"fmt"
	"sync"

	"github.com/go-errors/errors"
)

// Parallelize executes the given tasks parallelly and returns the result of execution.
// The resultset returned is in the same order as the tasks were given. So developers
// can rest assured that the execution order is not altered.
//
// It also panic friendly and tries best to deal with any panics occuring due to
// foreign function calls.
//
//
func Parallelize(functions []func() interface{}) []interface{} {
	max := len(functions)
	resultSet := make([]interface{}, max)

	wg := sync.WaitGroup{}
	wg.Add(max)
	for k, f := range functions {
		go func(f func() interface{}, offset int) {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Attention!!!   Panic Occurred !!!")
					fmt.Println("Handled Gracefully !!!")
					fmt.Println(errors.Wrap(r, 2).ErrorStack())
				}
			}()
			resultSet[offset] = f()

		}(f, k)

	}
	wg.Wait()
	return resultSet
}

// ParallelizeThrottled is same as Parallelize with an extra feature to limit number of parallel
// calls. It is useful, when the functions make calls to external machines or devices and calling
// several hundred of parallel goroutines could crash the external system.
//
// The function accepts a throttling factor as the second argument.
// A factor of 2 means, it will make atmost 2 parallel calls.
// A factor of 5 means, it will make atmost 5 parallel calls.
// A factor of 0 means, it will make all calls in parallel.
//
// It also panic friendly and tries best to deal with any panics occuring due to
// foreign function calls.
//
// Ordering of the resultset is maintained. So, developers can rest assured that the
// ouput order is same as input order.
//
func ParallelizeThrottled(functions []func() interface{}, factor int) []interface{} {

	if factor <= 0 {
		return Parallelize(functions)
	}
	max := len(functions)
	resultSet := make([]interface{}, max)

	var completed int //keep track of completed tasks.

	for i := 0; i < max; i = i + factor {
		var toProcess int
		if max-completed >= factor {
			toProcess = factor
		} else {
			toProcess = max - completed
		}
		wg := sync.WaitGroup{}
		wg.Add(toProcess)

		current := i
		for current < i+toProcess {
			go func(f func() interface{}, offset int) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						fmt.Println("Attention!!!   Panic Occurred !!!")
						fmt.Println("Handled Gracefully !!!")
						fmt.Println(errors.Wrap(r, 2).ErrorStack())
					}
				}()
				resultSet[offset] = f()

			}(functions[current], current)
			current++
		}
		wg.Wait()
		completed = completed + toProcess
	}
	return resultSet
}
