package arrayutils

import (
	"math"
	"runtime"
	"utils/typeutils"
)

func ExistsAt[T any](xs []T, compare typeutils.Compare[T]) (int, bool) {
	n := len(xs)
	if n == 0 {
		return 0, false
	}
	size := int(math.Ceil(float64(n) / float64(runtime.NumCPU())))
	slices := int(math.Ceil(float64(n) / float64(size)))
	result := make(chan *int, slices)
	var j int
	// ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < n; i += size {
		j += size
		if j > n {
			j = n
		}
		go func(x, y int) {
			cpart := xs[x:y]
			for i, n := 0, len(cpart); i < n; i++ {
				e := cpart[i]
				if compare(e) == typeutils.Equal {
					exitsAt := x + i
					result <- &exitsAt
					return
				}
			}
			result <- nil
		}(i, j)
	}
	// defer cancel()
	for i := 0; i < slices; i++ {
		v := <-result
		if v != nil {
			return *v, true
		}
	}
	return 0, false
}

func RemoveAt[V any](xs []V, i int) []V {
	return append(xs[:i], xs[i+1:]...)
}
