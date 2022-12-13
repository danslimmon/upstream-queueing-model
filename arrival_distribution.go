package main

import (
	"math/rand"
)

type ArrivalDistribution struct {
	avgArrivalInterval int
}

func (dist *ArrivalDistribution) ArrivalInterval() int {
	return int(rand.ExpFloat64() * float64(dist.avgArrivalInterval))
}

// avgReqPerSecond is the requests per second received by the upstream.
//
// Just a constant-interval drip.
func NewArrivalDistribution(avgReqPerSecond float64) *ArrivalDistribution {
	dist := new(ArrivalDistribution)
	dist.avgArrivalInterval = int(1000000.0 / avgReqPerSecond)
	return dist
}
