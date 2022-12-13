package main

import (
	"math/rand"
)

type LatencyBucket struct {
	Index        float64
	LatencyMicro int
}

type LatencyDistribution struct {
	buckets []LatencyBucket
}

func (dist *LatencyDistribution) ProcessingInterval() int {
	// pick a bucket to choose the latency from
	var min, max int
	bucketR := rand.Float64()
	for i, p := range dist.buckets {
		if bucketR < p.Index {
			min = dist.buckets[i-1].LatencyMicro
			max = p.LatencyMicro
			break
		}
	}

	v := min + rand.Intn(max-min)
	return v
}

func NewLatencyDistribution(buckets []LatencyBucket) (*LatencyDistribution, error) {
	return &LatencyDistribution{
		buckets: buckets,
	}, nil
}
