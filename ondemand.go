package traefik_ondemand_plugin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

const defaultTimeoutSeconds = 60

// Config the plugin configuration
type Config struct {
	Name       string
	ServiceUrl string
	Timeout    uint64
}

func CreateConfig() *Config {
	return &Config{
		Timeout: defaultTimeoutSeconds,
	}
}

type Ondemand struct {
	next              http.Handler
	name              string
	serviceUrl        string
	timeoutSeconds    uint64
	dockerServiceName string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.ServiceUrl) == 0 {
		return nil, fmt.Errorf("serviceUrl cannot be null")
	}

	if len(config.Name) == 0 {
		return nil, fmt.Errorf("name cannot be null")
	}

	return &Ondemand{
		next:              next,
		name:              name,
		serviceUrl:        config.ServiceUrl,
		dockerServiceName: config.Name,
		timeoutSeconds:    config.Timeout,
	}, nil
}

func (e *Ondemand) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	url := fmt.Sprintf("%s/?name=%s&timeout=%d", e.serviceUrl, e.dockerServiceName, e.timeoutSeconds)
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
