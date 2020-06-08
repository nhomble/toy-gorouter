package main

import "testing"

func TestRoundRobinStrategy_Next_Exhausted(t *testing.T) {
	rr := roundRobinStrategy{0}
	s := routeDestination{
		healthy: false,
	}
	pool := destinationPool{
		destinations: []*routeDestination{
			&s,
		},
	}
	out := rr.Next(&pool)
	if out != nil {
		t.Fatalf("Should have exhausted")
	}
}

func TestRoundRobinStrategy_Next(t *testing.T) {
	rr := roundRobinStrategy{0}
	bad := routeDestination{
		healthy: false,
	}
	good := routeDestination{
		healthy: true,
	}
	pool := destinationPool{
		destinations: []*routeDestination{
			&bad,
			&good,
		},
	}
	out := rr.Next(&pool)
	if out != &good {
		t.Fatalf("Should have returned the healthy routeDestination")
	}
}

func TestRoundRobinStrategy_Next_NilPool(t *testing.T) {
	rr := roundRobinStrategy{0}
	out := rr.Next(nil)
	if out != nil {
		t.Fatalf("Should have been nil")
	}
}
