package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLatencyDistribution(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dist, err := NewLatencyDistribution([]LatencyBucket{
		LatencyBucket{0.0, 10000},
		LatencyBucket{0.5, 20000},
		LatencyBucket{1.0, 40000},
	})
	assert.Nil(err)

	var sum int
	for i := 0; i < 10000; i++ {
		sum += dist.ProcessingInterval()
	}
	rsltAvg := sum / 10000
	expAvg := (15000 + 30000) / 2
	p := float64(rsltAvg) / float64(expAvg)
	// within a tenth
	assert.True(p > 0.9 && p < 1.11)
}
