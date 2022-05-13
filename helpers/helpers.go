package helpers

import (
	"math"
	"math/rand"
	"time"
)

const (
	BackoffFactorTwo   = float64(2)
	BackoffFactorFour  = float64(4)
	BackoffFactorEight = float64(8)
)

func GetExponentialBackoff(backoff float64, factor float64) time.Duration {
	min := 100 * time.Millisecond
	max := 30 * time.Second
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
