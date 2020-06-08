package main

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestConfigLoad(t *testing.T) {
	base, _ := os.Getwd()
	path := filepath.FromSlash("../resources/basic.yml")
	config, _ := readConfig(filepath.Join(base, path))
	if config.Port != "8080" {
		t.Fatalf("config.Port != 8080, instead it was %s", config.Port)
	}
	if len(config.Backends) != 2 {
		t.Fatalf("len(config.Backends) != 2, instead it was %d", len(config.Backends))
	}
	if config.Backends[0] != "http://localhost:8081" {
		t.Fatalf("config.Backends[0] != http://localhost:8081, instead it was %s", config.Backends[0])
	}
	if config.Backends[1] != "http://localhost:8082" {
		t.Fatalf("config.Backends[1] != http://localhost:8082, instead it was %s", config.Backends[1])
	}
}

func TestBadConfigLoad(t *testing.T) {
	base, _ := os.Getwd()
	path := filepath.FromSlash("../resources/blank.yml")
	_, err := readConfig(filepath.Join(base, path))
	if err == nil || err.Error() != "Must define a port" {
		log.Fatalf("Expected error!")
	}
}
