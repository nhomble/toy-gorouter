package main

import (
	"sync/atomic"
)

type roundRobinStrategy struct {
	current uint64
}

func (strat *roundRobinStrategy) increment(max uint64) {
	if max > 0 {
		atomic.AddUint64(&strat.current, uint64(1))
		atomic.StoreUint64(&strat.current, atomic.LoadUint64(&strat.current)%max)
	}
}

func (strat *roundRobinStrategy) Next(pool *destinationPool) *routeDestination {
	if pool == nil {
		return nil
	}
	strat.increment(uint64(len(pool.destinations)))
	s := pool.destinations[strat.current]
	for exhausted := 0; exhausted < 2*len(pool.destinations); exhausted++ {
		logDebug("roundRobin index=%d healthy=%v", strat.current, s.CheckHealth())
		if s.CheckHealth() {
			return s
		}
		strat.increment(uint64(len(pool.destinations)))
		s = pool.destinations[strat.current]
	}
	return nil
}
