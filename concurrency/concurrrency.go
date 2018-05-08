package concurrency

func Parallelize(functions []func() interface{}) []interface{} {

	ch := make(chan interface{})
	for _, f := range functions {
		go func(f func() interface{}) {
			ch <- f()
		}(f)
	}

	max := len(functions)
	var count int
	i := make([]interface{}, 0)
	for {
		if count == max {
			break
		}
		count++
		i = append(i, <-ch)
	}
	return i
}
