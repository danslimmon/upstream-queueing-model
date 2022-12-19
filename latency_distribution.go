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
    /*
    // Used to calculate average latency
    sum int
    n int
    */
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
    /*
    dist.sum += v
    dist.n += 1
    */
	return v
}

/*
// Average returns the average of all values returned from ProcessingInterval so far.
//
// Returns -1 if called before any calls to ProcessingInterval.
func (dist *LatencyDistribution) Average() int {
    if dist.n == 0 {
        return -1
    }
    return int(float64(dist.sum) / float64(dist.n) + 0.5)
}

func (dist *LatencyDistribution) Variance() 
*/

func NewLatencyDistribution(buckets []LatencyBucket) (*LatencyDistribution, error) {
	return &LatencyDistribution{
		buckets: buckets,
	}, nil
}
