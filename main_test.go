package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventStream(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	{
		stream := new(EventStream)
		exp := Event{
			Time: 0,
			Type: "arrival",
		}
		stream.Schedule(exp)
		got := stream.Next()
		assert.Equal(exp, got)
	}

	{
		// two events scheduled in time order
		stream := new(EventStream)
		exp0 := Event{
			Time: 0,
			Type: "arrive",
		}
		exp1 := Event{
			Time: 5,
			Type: "finish",
		}
		stream.Schedule(exp0)
		stream.Schedule(exp1)
		got := stream.Next()
		assert.Equal(exp0, got)
		got = stream.Next()
		assert.Equal(exp1, got)
	}

	{
		// two events scheduled in reverse time order
		stream := new(EventStream)
		exp0 := Event{
			Time: 0,
			Type: "arrive",
		}
		exp1 := Event{
			Time: 5,
			Type: "finish",
		}
		stream.Schedule(exp1)
		stream.Schedule(exp0)
		got := stream.Next()
		assert.Equal(exp0, got)
		got = stream.Next()
		assert.Equal(exp1, got)
	}

	{
		// two events with same time (should be popped in order scheduled)
		stream := new(EventStream)
		exp0 := Event{
			Time: 5,
			Type: "arrive",
		}
		exp1 := Event{
			Time: 5,
			Type: "finish",
		}
		stream.Schedule(exp0)
		stream.Schedule(exp1)
		got := stream.Next()
		assert.Equal(exp0, got)
		got = stream.Next()
		assert.Equal(exp1, got)
	}

	{
		// pop on empty EventStream
		stream := new(EventStream)
		assert.Panics(func() {
			stream.Next()
		})
	}
}

func TestUpstreamRun(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	{
		u := NewUpstream()
		in := SimInput{
			NRequests:          10,
			MaxInflight:        1,
			ArrivalInterval:    func() int { return 10 },
			ProcessingInterval: func() int { return 5 },
		}
		out := u.Run(in)
		assert.Equal(0, out.QueuedTime)
		assert.Equal(10*5, out.TotalTime)
		assert.Equal(10*5, out.ProcessingTime)
		assert.Equal(95, out.CapacityTime)
	}

	{
		u := NewUpstream()
		in := SimInput{
			NRequests:          10,
			MaxInflight:        1,
			ArrivalInterval:    func() int { return 10 },
			ProcessingInterval: func() int { return 15 },
		}
		out := u.Run(in)

		// 5 second intervals:
		//
		// Queued:    0 0 1   0 1 1   1 1 2   1 2 2   2 2 3   2 3 3   3 3 4   3 4 4   4 4 5   4 5 5
		// Inflight:  1 1 1   1 1 1   1 1 1   1 1 1   1 1 1   1 1 1   1 1 1   1 1 1   1 1 1   1 1 1
		// Arrival:  ^   ^     ^   ^     ^     ^   ^     ^     ^   ^     ^     ^   ^     ^     ^   ^
		// Finish:         ^       ^       ^       ^       ^       ^       ^       ^       ^       ^
		// Dequeue:        ^       ^       ^       ^       ^       ^       ^       ^       ^
		assert.Equal(5*((0+0+1)+(0+1+1)+(1+1+2)+(1+2+2)+(2+2+3)+(2+3+3)+(3+3+4)+(3+4+4)+(4+4+5)+(4+5+5)), out.QueuedTime)
		assert.Equal(10*15+out.QueuedTime, out.TotalTime)
		assert.Equal(10*15, out.ProcessingTime)
		assert.Equal(out.ProcessingTime, out.CapacityTime)
	}
}
