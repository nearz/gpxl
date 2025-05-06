package concurrent

import (
	"runtime"
	"sync"
)

func Rows(ylen, ymin int, process func(start, end int)) {
	n := runtime.NumCPU()
	chunk := (ylen + n - 1) / n
	if n <= 1 || chunk <= n {
		process(0, ylen)
	} else {
		var wg sync.WaitGroup
		for sy := ymin; sy < ylen; sy += chunk {
			wg.Add(1)
			ey := sy + chunk
			if ey > ylen {
				ey = ylen
			}
			go func() {
				defer wg.Done()
				process(sy, ey)
			}()
		}
		wg.Wait()
	}
}
