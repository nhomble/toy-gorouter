package main

import (
	"net/http"
	"time"
)

type httpHealthcheck struct {
	interval int
	endpoint string
	client   *http.Client
}

func newHttpHealthcheck(interval int, timeout int, endpoint string) httpHealthcheck {
	client := &http.Client{
		Timeout: time.Second * time.Duration(timeout),
	}
	return httpHealthcheck{
		interval: interval,
		endpoint: endpoint,
		client:   client,
	}
}

func (h httpHealthcheck) Loop(p *destinationPool) {
	tick := time.NewTicker(time.Second * time.Duration(h.interval))
	for {
		select {
		case <-tick.C:
			logDebug("Starting health check loop - http checks")
			for _, ele := range p.destinations {
				health := h.Check(ele)
				ele.SetHealth(health)
				if health {
					logInfo("http health result url=%s%s result=%t", ele.url.String(), h.endpoint, health)
				} else {
					logWarn("http health result url=%s%s result=%t", ele.url.String(), h.endpoint, health)
				}

			}
		}
	}
}

func (h httpHealthcheck) Check(r *routeDestination) bool {
	resp, err := h.client.Get(r.url.String() + h.endpoint)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}
