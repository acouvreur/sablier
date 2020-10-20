package service

import (
	"time"
)

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
