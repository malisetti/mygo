package utils

import (
	"math"
	"runtime"
)

type CompareFunc[T any] func(x T) bool

func ExistsAt[T any](xs []T, compare CompareFunc[T]) (int, bool) {
	n := len(xs)
	if n == 0 {
		return 0, false
	}
	size := int(math.Ceil(float64(n) / float64(runtime.NumCPU())))
	slices := int(math.Ceil(float64(n) / float64(size)))
	result := make(chan *int, slices)
	var j int
	for i := 0; i < n; i += size {
		j += size
		if j > n {
			j = n
		}
		go func(x, y int) {
			cpart := xs[x:y]
			for i := 0; i < len(cpart); i++ {
				at := i
				e := cpart[at]
				if compare(e) {
					exitsAt := x + at
					result <- &exitsAt
					return
				}
			}
			result <- nil
		}(i, j)
	}
	for i := 0; i < slices; i++ {
		v := <-result
		if v != nil {
			return *v, true
		}
	}
	return 0, false
}
