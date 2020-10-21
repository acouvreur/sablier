package traefik_ondemand_plugin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

const defaultTimeoutSeconds = 10

// Config the plugin configuration
type Config struct {
	DockerServiceName string
	ServiceUrl        string
	TimeoutSeconds    uint64
}

func CreateConfig() *Config {
	return &Config{
		TimeoutSeconds: defaultTimeoutSeconds,
	}
}

type Ondemand struct {
	next              http.Handler
	name              string
	ServiceUrl        string
	TimeoutSeconds    uint64
	DockerServiceName string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.ServiceUrl) == 0 {
		return nil, fmt.Errorf("ServiceUrl cannot be null")
	}

	return &Ondemand{
		next:              next,
		name:              name,
		ServiceUrl:        config.ServiceUrl,
		DockerServiceName: config.DockerServiceName,
		TimeoutSeconds:    config.TimeoutSeconds,
	}, nil
}

func (e *Ondemand) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("%s/?name=%s&timeout=%d", e.ServiceUrl, e.DockerServiceName, e.TimeoutSeconds)
	resp, err := http.Get(url)
	if err != nil {
		println("Could not contact", url)
		e.next.ServeHTTP(rw, req)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println("Could not parse service response")
		e.next.ServeHTTP(rw, req)
		return
	}

	fmt.Printf("%s\n", body)
	bodystr := string(body)

	if bodystr == "started" {
		// Service started forward request
		e.next.ServeHTTP(rw, req)
	} else if bodystr == "starting" {
		// Service starting, notify client
		rw.Write([]byte("Service is starting..."))
	} else {
		// Error
		rw.Write(body)
	}
}
