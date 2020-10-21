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

	queryParams := r.URL.Query()
	fmt.Printf("%+v\n", queryParams)
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
	cli, err := client.NewEnvClient()
	if err != nil {
		fmt.Fprintf(w, "%+v", "Could not connect to docker API")
	}
	// 1. Check if service is up
	// 2.
	// IS DOWN
	// 2.1 Start the service if down (async)
	// 2.2 Set timeout
	// IS UP
	// 2.1 Reset timeout
	// 3. Response

	service := GetOrCreateService(serviceName, uint64(serviceTimeout))
	err = service.HandleServiceState(cli)
	if err != nil {
		fmt.Printf("Error: %+v ", err)
		fmt.Fprintf(w, "%+v", err)
	}
	fmt.Printf("%+v\n", service)
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
func (service *Service) HandleServiceState(cli *client.Client) error {
	if service.isUp() == true {
		fmt.Printf("- Service %v is up\n", service.name)
		service.timeout = service.initialTimeout

	} else if service.isDown() {
		fmt.Printf("- Service %v is down\n", service.name)
		service.start()
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
	filterOPt := opts.NewFilterOpt()
	listOpts := types.ServiceListOptions{
		Filters: filterOPt.Value(),
	}
	services, err := client.ServiceList(ctx, listOpts)

	if err != nil {
		return err
	}

	dockerService, err := findService(services, service.name)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", dockerService)

	dockerService.Spec.Mode.Replicated = &swarm.ReplicatedService{
		Replicas: create(oneReplica),
	}

	swarmCluster, err := client.SwarmInspect(ctx)
	if err != nil {
		return err
	}

	client.ServiceUpdate(ctx, dockerService.ID, swarmCluster.ClusterInfo.Version, dockerService.Spec, types.ServiceUpdateOptions{})

	status := UP
	service.status = status
	return nil
}

func findService(services []swarm.Service, name string) (*swarm.Service, error) {
	for _, service := range services {
		if name == service.Spec.Name {
			return &service, nil
		}
	}
	return &swarm.Service{}, fmt.Errorf("Could not find service %s", name)
}

func (service *Service) start() {
	// start service in docker
	service.timeout = service.initialTimeout
	go service.stopAfterTimeout()
}

func (service *Service) stopAfterTimeout() {
	for service.timeout > 0 {
		time.Sleep(100 * time.Millisecond)
	}
	// TODO :: stop the service
}

func (service *Service) setServiceReplicas(client *client.Client) error {
	ctx := context.Background()
	filterOPt := opts.NewFilterOpt()
	listOpts := types.ServiceListOptions{
		Filters: filterOPt.Value(),
	}
	services, err := client.ServiceList(ctx, listOpts)
	if err != nil {
		return fmt.Errorf("Error: %+v", err)
	}

	fmt.Printf("%+v", services)

	return nil

}

func create(x uint64) *uint64 {
	return &x
}
