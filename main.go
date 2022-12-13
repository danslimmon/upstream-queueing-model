package main

import (
	"fmt"
)

type Event struct {
	Time int
	Type string
}

type EventStream struct {
	events []Event
}

func (stream *EventStream) Schedule(ev Event) {
	initialLen := len(stream.events)
	for i := range stream.events {
		if stream.events[i].Time > ev.Time {
			stream.events = append(stream.events[0:i+1], stream.events[i:len(stream.events)]...)
			stream.events[i] = ev
			return
		}
	}
	if len(stream.events) == initialLen {
		stream.events = append(stream.events, ev)
	}
}

func (stream *EventStream) Next() Event {
	// Will panic if len(events) == 0. That's okay. There should always be at least one more arrival
	// event. If there isn't, we have a bug.
	ev := stream.events[0]
	stream.events = stream.events[1:len(stream.events)]
	return ev
}

type SimInput struct {
	NRequests          int
	MaxInflight        int
    // returns the number of microseconds before the next arrival of a request
	ArrivalInterval    func() int
    // returns the number of microseconds that a request will take to serve
	ProcessingInterval func() int
}

// times in microseconds
type SimOutput struct {
	QueuedTime     int
	TotalTime      int
	ProcessingTime int
	CapacityTime   int
}

type Upstream struct {
	inflight    int
	queued      int
	eventStream *EventStream

	// stats, in microseconds
	queuedTime int
	totalTime  int
}

func (u *Upstream) Run(in SimInput) SimOutput {
	i := 0
	prevTime := 0
	var ev Event
	for i < in.NRequests {
		ev = u.eventStream.Next()

		// update stats
		u.queuedTime += u.queued * (ev.Time - prevTime)
		u.totalTime += (u.queued + u.inflight) * (ev.Time - prevTime)

		switch ev.Type {
		case "arrive":
			// update counters
			if u.inflight >= in.MaxInflight {
				u.queued++
			} else {
				u.inflight++
				// schedule this request finishing
				u.eventStream.Schedule(Event{
					Time: ev.Time + in.ProcessingInterval(),
					Type: "finish",
				})
			}

			// schedule next arrival
			u.eventStream.Schedule(Event{
				Time: ev.Time + in.ArrivalInterval(),
				Type: "arrive",
			})
		case "finish":
			// update counters
			u.inflight--
			if u.queued > 0 {
				u.inflight++
				u.queued--

				// schedule the dequeued request finishing
				u.eventStream.Schedule(Event{
					Time: ev.Time + in.ProcessingInterval(),
					Type: "finish",
				})
			}

			// deal with loop
			i++
		}

		prevTime = ev.Time
	}

	return SimOutput{
		QueuedTime:     u.queuedTime,
		TotalTime:      u.totalTime,
		ProcessingTime: (u.totalTime - u.queuedTime),
		CapacityTime:   in.MaxInflight * ev.Time,
	}
}

func NewUpstream() *Upstream {
	stream := new(EventStream)
	stream.Schedule(Event{
		Time: 0,
		Type: "arrive",
	})
	return &Upstream{
		eventStream: stream,
	}
}

func main() {
    // percentiles come from system telemetry
	latencyDist, err := NewLatencyDistribution([]LatencyBucket{
		LatencyBucket{0.0, 10000},
		LatencyBucket{0.5, 75000},
		LatencyBucket{0.75, 192000},
		LatencyBucket{0.90, 601000},
		LatencyBucket{0.95, 1260000},
		LatencyBucket{0.99, 3410000},
		// actual max is much higher than this, but it's a long tail. definitely not uniformly
		// distributed between P99 and P100. hard to get more detail out of telemetry.
		LatencyBucket{1.0, 6420000},
	})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("queue_time,total_time,reqs_per_sec,pct_time_queued,pct_capacity_used,avg_secs_queued\n")
	for reqsPerSec := 14.0; reqsPerSec < 39.0; reqsPerSec += .05 {
		u := NewUpstream()
		arrivalDist := NewArrivalDistribution(reqsPerSec)

		in := SimInput{
			NRequests:          1000000,
			MaxInflight:        12,
			ProcessingInterval: latencyDist.ProcessingInterval,
			ArrivalInterval:    arrivalDist.ArrivalInterval,
		}
		out := u.Run(in)

		// queue_time
        //
        // Total amount of time spent by requests in upstream-local queue
		fmt.Printf("%f,", float64(out.QueuedTime)/1000000.0)

		// total_time
        //
        // Total time spent by requests in the system, queued or not
		fmt.Printf("%f,", float64(out.TotalTime)/1000000.0)

		// reqs_per_sec
        //
        // Average number of requests per second sent to the upstrea
		fmt.Printf("%f,", reqsPerSec)

		// pct_time_queued
        //
        // Percentage of request-seconds spent in the upstream-local queue
		fmt.Printf("%f,", float64(out.QueuedTime)/float64(out.TotalTime)*100.0)

		// pct_capacity_used
        //
        // Percentage of request-seconds spent  actively processing (i.e. not queued)
		fmt.Printf("%f,", float64(out.ProcessingTime)/float64(out.CapacityTime)*100.0)

		// avg_secs_queued
        //
        // The average number of seconds spent in the upstream-local queue by a request.
		fmt.Printf("%f\n", float64(out.QueuedTime)/float64(in.NRequests)/1000000.0)
	}
}
