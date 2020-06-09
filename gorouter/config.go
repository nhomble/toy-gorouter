package main

import (
	"errors"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

type RouterConfig struct {
	Port     string
	Backends []string `yaml:",flow"`
	Balancer BalancerConfig
	Health   HealthConfig
}

type BalancerConfig struct {
	Type string
}

type HealthConfig struct {
	Type     string
	Interval int
	Timeout  int
	Endpoint string
}

func readConfig(path string) (RouterConfig, error) {
	file, err := ioutil.ReadFile(path)
	data := RouterConfig{}

	if err != nil {
		return data, err
	}

	err = yaml.Unmarshal(file, &data)
	if len(data.Port) == 0 {
		return data, errors.New("Must define a port")
	}

	return data, err
}
