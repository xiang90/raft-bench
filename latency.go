package main

var (
	latency [100000][2]int64
)

func avgLatency() int64 {
	var (
		avg int64
		n   int
		l   [2]int64
	)

	for n, l = range latency {
		// do not count the warm up ones
		if n < 3 {
			continue
		}
		if l[0] == 0 {
			break
		}
		avg += l[1] - l[0]
	}
	return avg / int64(n-2)
}
