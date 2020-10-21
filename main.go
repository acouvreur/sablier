package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/swarm"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/opts"
)

const oneReplica = uint64(1)
const zeroReplica = uint64(0)

func handleRequests(w http.ResponseWriter, r *http.Request) {
	serviceName, serviceTimeout := parseParams(w, r)
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Fprintf(w, "%+v", "Could not connect to docker API")
	}
	service := GetOrCreateService(serviceName, serviceTimeout)
	err = service.HandleServiceState(cli)
	if err != nil {
		fmt.Printf("Error: %+v\n ", err)
		fmt.Fprintf(w, "%+v", err)
	}
	fmt.Printf("Service after query: %+v\n", service)
	fmt.Fprintf(w, "%+v", service)

}

func parseParams(w http.ResponseWriter, r *http.Request) (string, uint64) {
	queryParams := r.URL.Query()

	if queryParams["name"] == nil {
		http.Error(w, "name is required", 400)
		fmt.Fprintf(w, "%+v", "name is required")
	}
	serviceName := string(queryParams["name"][0])
	if queryParams["timeout"] == nil {
		http.Error(w, "timeout is required", 400)
		fmt.Fprintf(w, "%+v", "name is required")
	}

	serviceTimeout, err := strconv.Atoi(queryParams["timeout"][0])
	if err != nil {
		fmt.Fprintf(w, "%+v", "timeout must be an integer.")
	}
	return serviceName, uint64(serviceTimeout)
}

func main() {
	http.HandleFunc("/", handleRequests)
	log.Fatal(http.ListenAndServe(":10000", nil))
}

// ===  Other file ===

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
func (service *Service) HandleServiceState(cli *client.Client) error {
	if service.isUp() == true {
		fmt.Printf("- Service %v is up\n", service.name)
		service.timeout = service.initialTimeout
		go service.stopAfterTimeout(cli)
	} else if service.isDown() {
		fmt.Printf("- Service %v is down\n", service.name)
		service.start(cli)
	} else {
		fmt.Printf("- Service %v status is unknown\n", service.name)
		err := service.setServiceStateFromDocker(cli)
		if err != nil {
			return err
		}
		service.HandleServiceState(cli)
	}
	return nil
}

func (service *Service) isUp() bool {
	return service.status == UP
}

func (service *Service) isDown() bool {
	return service.status == DOWN
}

func (service *Service) setServiceStateFromDocker(client *client.Client) error {
	ctx := context.Background()
	dockerService, err := service.getDockerService(ctx, client)

	if err != nil {
		return err
	}

	status := UP
	fmt.Printf("replicas %d\n", dockerService.Spec.Mode.Replicated.Replicas)
	if *dockerService.Spec.Mode.Replicated.Replicas == zeroReplica {
		status = DOWN
	}
	service.status = status
	return nil
}

func (service *Service) start(client *client.Client) {
	fmt.Printf("Starting service %s\n", service.name)
	service.setServiceReplicas(client, 1)
	service.timeout = service.initialTimeout
	go service.stopAfterTimeout(client)
}

func (service *Service) stopAfterTimeout(client *client.Client) {
	for service.timeout > 0 {
		time.Sleep(1 * time.Second)
		service.timeout--
	}
	fmt.Printf("Stopping service %s\n", service.name)
	service.setServiceReplicas(client, 0)
}

func (service *Service) setServiceReplicas(client *client.Client, replicas uint64) error {
	ctx := context.Background()
	dockerService, err := service.getDockerService(ctx, client)
	if err != nil {
		return err
	}
	dockerService.Spec.Mode.Replicated = &swarm.ReplicatedService{
		Replicas: create(replicas),
	}
	client.ServiceUpdate(ctx, dockerService.ID, dockerService.Meta.Version, dockerService.Spec, types.ServiceUpdateOptions{})
	return nil

}

func (service *Service) getDockerService(ctx context.Context, client *client.Client) (*swarm.Service, error) {
	filterOPt := opts.NewFilterOpt()
	listOpts := types.ServiceListOptions{
		Filters: filterOPt.Value(),
	}
	services, err := client.ServiceList(ctx, listOpts)

	if err != nil {
		return nil, err
	}

	dockerService, err := findService(services, service.name)

	if err != nil {
		return nil, err
	}

	return dockerService, nil
}

func findService(services []swarm.Service, name string) (*swarm.Service, error) {
	for _, service := range services {
		if name == service.Spec.Name {
			return &service, nil
		}
	}
	return &swarm.Service{}, fmt.Errorf("Could not find service %s", name)
}

func create(x uint64) *uint64 {
	return &x
}
