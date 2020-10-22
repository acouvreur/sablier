package traefik_ondemand_plugin

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	typeName = "Ondemand"
)

const defaultTimeoutSeconds = 60

// Net client is a custom client to timeout after 2 seconds if the service is not ready
var netClient = &http.Client{
	Timeout: time.Second * 2,
}

// Config the plugin configuration
type Config struct {
	Name       string
	ServiceUrl string
	Timeout    uint64
}

// CreateConfig creates a config with its default values
func CreateConfig() *Config {
	return &Config{
		Timeout: defaultTimeoutSeconds,
	}
}

// Ondemand holds the request for the on demand service
type Ondemand struct {
	request string
	name    string
	next    http.Handler
}

func buildRequest(url string, name string, timeout uint64) (string, error) {
	// TODO: Check url validity
	request := fmt.Sprintf("%s?name=%s&timeout=%d", url, name, timeout)
	return request, nil
}

// New function creates the configuration
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if len(config.ServiceUrl) == 0 {
		return nil, fmt.Errorf("serviceUrl cannot be null")
	}

	if len(config.Name) == 0 {
		return nil, fmt.Errorf("name cannot be null")
	}

	request, err := buildRequest(config.ServiceUrl, config.Name, config.Timeout)

	if err != nil {
		return nil, fmt.Errorf("error while building request")
	}

	return &Ondemand{
		next:    next,
		name:    name,
		request: request,
	}, nil
}

// ServeHTTP retrieve the service status
func (e *Ondemand) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	status, err := getServiceStatus(e.request)

	println(status, err == nil)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}

	if status == "started" {
		println("Started !")
		// Service started forward request
		e.next.ServeHTTP(rw, req)

	} else if status == "starting" {
		println("Starting !")
		// Service starting, notify client
		rw.WriteHeader(http.StatusAccepted)
		rw.Write([]byte("Service is starting..."))
	} else {
		println("Error :() !")
		// Error
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("Unexpected status answer from ondemand service"))
	}
}

func getServiceStatus(request string) (string, error) {

	// This request wakes up the service if he's scaled to 0
	println(request)
	resp, err := netClient.Get(request)
	if err != nil {
		return "error", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "parsing error", err
	}

	return strings.TrimSuffix(string(body), "\n"), nil
}
