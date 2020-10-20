package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handleRequests(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	fmt.Printf("%+v", queryParams)
	if queryParams["name"] == nil {
		http.Error(w, "name is required", 400)
		fmt.Fprintf(w, "%+v", "name is required")
	}
	if queryParams["timeout"] == nil {
		http.Error(w, "timeout is required", 400)
		fmt.Fprintf(w, "%+v", "name is required")
	}
	// 1. Check if service is up
	// 2.
	// IS DOWN
	// 2.1 Start the service if down (async)
	// 2.2 Set timeout
	// IS UP
	// 2.1 Reset timeout
	// 3. Response

	service := GetOrCreateService("test", 10)
	// service.HandleServiceState(
	fmt.Fprintf(w, "%+v", service)

}

func main() {
	http.HandleFunc("/", handleRequests)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

/// Other file

// Status is the service status
type Status string

const (
	UP      Status = "up"
	DOWN    Status = "down"
	UNKNOWN Status = "unknown"
)

// Service holds all information related to a service
type Service struct {
	name           string
	timeout        uint64
	initialTimeout uint64
	status         Status
}

var services = map[string]*Service{}

// GetOrCreateService return an existing service or create one
func GetOrCreateService(name string, timeout uint64) *Service {
	if services[name] != nil {
		return services[name]
	}
	service := &Service{name, timeout, timeout, UNKNOWN}

	services[name] = service
	return service
}

// HandleServiceState up the service if down or set timeout for downing the service
func (service *Service) HandleServiceState() {
	if service.isUp() == true {
		service.timeout = service.initialTimeout
		go service.stopAfterTimeout()

	} else if service.isDown() {
		service.start()
	} else {
		service.setServiceStateFromDocker()
		service.HandleServiceState()
	}
}

func (service *Service) isUp() bool {
	return service.status == UP
}

func (service *Service) isDown() bool {
	return service.status == DOWN
}

func (service *Service) setServiceStateFromDocker() {
	// set status form docker
	status := UNKNOWN
	service.status = status
}

func (service *Service) start() {
	// start service in docker
	service.timeout = service.initialTimeout
}

func (service *Service) stopAfterTimeout() {
	for service.timeout > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	println("OVER MOTHER FUCKER")
}
