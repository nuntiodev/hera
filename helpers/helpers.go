package helpers

import (
	"math"
	"math/rand"
	"time"
)

func GetExponentialBackoff(backoff float64) time.Duration {
	min := 100 * time.Millisecond
	max := 10 * time.Second
	factor := float64(2)
	minf := float64(min)
	durf := minf * math.Pow(factor, backoff)
	durf = rand.Float64()*(durf-minf) + minf
	dur := time.Duration(durf)
	//keep within bounds
	if dur < min {
		return min
	} else if dur > max {
		return max
	} else {
		return dur
	}
}
