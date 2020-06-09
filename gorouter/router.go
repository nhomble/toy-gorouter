package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

const (
	RETRY = iota
)

type routeStrategy interface {
	Next(pool *destinationPool) *routeDestination
}

type healthCheck interface {
	Loop(pool *destinationPool)
	Check(dest *routeDestination) bool
}

// A Server is a backend resource this router can route to. A routeDestination is defined in the configuration file and can be
// unhealthy
type routeDestination struct {
	url     *url.URL
	healthy bool
	proxy   *httputil.ReverseProxy
	lock    sync.RWMutex
}

func (s *routeDestination) SetHealth(b bool) {
	s.lock.Lock()
	s.healthy = b
	s.lock.Unlock()
}

func (s *routeDestination) CheckHealth() (health bool) {
	s.lock.RLock()
	health = s.healthy
	s.lock.RUnlock()
	return
}

type destinationPool struct {
	destinations []*routeDestination // all known destinations
	strategy     routeStrategy
}

func (p *destinationPool) destinationForUrl(url *url.URL) (r *routeDestination) {
	for _, ele := range p.destinations {
		logDebug("Comparing %v <=> %v", ele.url, url)
		if ele.url.Host == url.Host {
			r = ele
			return
		}
	}
	return nil
}

type routingHandler struct {
	config RouterConfig
	pool   *destinationPool
	health healthCheck
}

func getRetryCount(r *http.Request) (retry int) {
	retry, ok := r.Context().Value(RETRY).(int)
	if !ok {
		retry = 0
	}
	return
}

// Make a routable destination for the server
// The complexity of the function comes from defining the error handler in the context of the routeDestination
func makeDestination(handler *routingHandler, s string) routeDestination {
	u, e := url.Parse(s)
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, e error) {
		logError("Error requesting to host=%s", e, r.Host)
		retries := getRetryCount(r)
		if retries < 3 {
			select {
			case <-time.After(10 * time.Millisecond):
				logDebug("Going to retry request with retryCount=%d", retries+1)
				ctx := context.WithValue(r.Context(), RETRY, retries+1) // increment retry count
				proxy.ServeHTTP(w, r.WithContext(ctx))
			}
			return
		}
		logWarn("Marking url=%s as UNHEALTHY", s)
		inError := handler.pool.destinationForUrl(u)
		inError.SetHealth(false)
		handler.ServeHTTP(w, r)
	}
	return routeDestination{
		url:     u,
		healthy: e == nil,
		proxy:   proxy,
		lock:    sync.RWMutex{},
	}
}

// Create a routing handler. This includes:
// 1. creating a pool of backend destinations
// 2. initializing balancing strategy
// 3. initializing health checks
func makeRoutingHandler(config RouterConfig) routingHandler {
	handler := routingHandler{}
	handler.config = config
	handler.pool = new(destinationPool)
	handler.pool.destinations = make([]*routeDestination, len(handler.config.Backends))
	for i, s := range handler.config.Backends {
		srv := makeDestination(&handler, s)
		handler.pool.destinations[i] = &srv
		logDebug("Registering for routeDestination=%s", s)
	}
	if config.Balancer.Type == "roundrobin" {
		handler.pool.strategy = &roundRobinStrategy{0}
	}
	if config.Health.Type == "tcp" {
		handler.health = tcpHealthcheck{
			timeout:  config.Health.Timeout,
			interval: config.Health.Interval,
		}
	} else if config.Health.Type == "http" {
		handler.health = newHttpHealthcheck(config.Health.Interval, config.Health.Timeout, config.Health.Endpoint)
	}

	go handler.health.Loop(handler.pool)
	return handler
}

// Core routing logic here. In this handler (web filter) we:
// 1. find next routeDestination
// 2. delegate request handling to the proxy handler defined
func (c *routingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := c.pool.strategy.Next(c.pool)
	if n != nil {
		n.proxy.ServeHTTP(w, r)
	} else {
		http.Error(w, "Service not available", http.StatusServiceUnavailable)
	}
}
