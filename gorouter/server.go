package main

import (
	"fmt"
	"log"
	"net/http"
)

// Start the router routeDestination
func start(config RouterConfig) {
	handler := makeRoutingHandler(config)
	server := http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: &handler,
	}

	defer func() {
		err := server.Close()
		if err != nil {
			log.Fatal(err)
		} else {
			logInfo("Stopping router")
		}
	}()
	logInfo("Router starting on port=%s\n", config.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
