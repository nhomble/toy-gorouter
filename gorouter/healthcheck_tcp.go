package main

import (
	"net"
	"time"
)

type tcpHealthcheck struct {
	interval int
	timeout  int
}

func (t tcpHealthcheck) Loop(p *destinationPool) {
	tick := time.NewTicker(time.Second * time.Duration(t.interval))
	for {
		select {
		case <-tick.C:
			logDebug("Starting health check loop - TCP checks")
			for _, ele := range p.destinations {
				logDebug("Checking health of url=%s", ele.url.Host)
				health := t.Check(ele)
				ele.SetHealth(health)
				if health {
					logInfo("TCP health result url=%s result=%t", ele.url.Host, health)
				} else {
					logWarn("TCP health result url=%s result=%t", ele.url.Host, health)
				}
			}
		}
	}
}

func (t tcpHealthcheck) Check(r *routeDestination) bool {
	timeout := time.Second * time.Duration(t.timeout)
	conn, err := net.DialTimeout("tcp", r.url.Host, timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
