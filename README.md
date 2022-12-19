# upstream-queueing-model

This computational model simulates the behavior of an "upstream" which is being fed requests at
random intervals by a load balancer, and taking a random amount of time to process them. Each
upstream has a maximum concurrency of 12 before queueing latency is incurred.

See [this blog
post](https://blog.danslimmon.com/2022/12/19/round-robin-load-balancing-latency-and-queues/) to
understand how this model can be used.
