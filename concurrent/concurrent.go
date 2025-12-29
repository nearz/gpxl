package concurrent

import (
	"runtime"
	"sync"
)

const minChunk = 16

func Rows(ylen, ymin int, process func(start, end int)) {
	n := runtime.NumCPU()
	rows := ylen - ymin
	if n <= 1 || rows < n*minChunk {
		process(ymin, ylen)
		return
	}

	chunk := (rows + n - 1) / n

	var wg sync.WaitGroup
	for sy := ymin; sy < ylen; sy += chunk {
		wg.Add(1)
		ey := sy + chunk
		if ey > ylen {
			ey = ylen
		}
		go func(start, end int) {
			defer wg.Done()
			process(start, end)
		}(sy, ey)
	}
	wg.Wait()
}
