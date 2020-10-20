package traefik_ondemand_plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type ServiceResponse struct {
	status string `json:status`
}

func (e *Ondemand) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	resp, err := http.Get(e.ServiceUrl)
	if err != nil {
		println("Could not contact", e.ServiceUrl)
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
	serviceResponse := &ServiceResponse{}
	err = json.Unmarshal(body, serviceResponse)
	if err != nil {
		fmt.Println("error:", err)
		e.next.ServeHTTP(rw, req)
		return
	}
	fmt.Printf("%+v\n", serviceResponse)
	e.next.ServeHTTP(rw, req)
}
