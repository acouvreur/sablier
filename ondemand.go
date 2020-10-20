package traefik_ondemand_plugin

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

const defaultTimeoutSeconds = 60

// Config the plugin configuration
type Config struct {
	ServiceUrl     string
	TimeoutSeconds uint64
}

func CreateConfig() *Config {
	return &Config{
		TimeoutSeconds: defaultTimeoutSeconds,
	}
}

type Ondemand struct {
	next           http.Handler
	name           string
	ServiceUrl     string
	TimeoutSeconds uint64
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.ServiceUrl) == 0 {
		return nil, fmt.Errorf("ServiceUrl cannot be null")
	}

	return &Ondemand{
		next:           next,
		name:           name,
		ServiceUrl:     config.ServiceUrl,
		TimeoutSeconds: config.TimeoutSeconds,
	}, nil
}

func (e *Ondemand) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("%+v\n", e.ServiceUrl)
	log.Println("plugin executed")

	e.next.ServeHTTP(rw, req)
}
