package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/docker/docker/api/types/swarm"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/opts"
)

const oneReplica = uint64(1)
const zeroReplica = uint64(0)

// Status is the service status
type Status string

const (
	UP       Status = "up"
	DOWN     Status = "down"
	STARTING Status = "starting"
	UNKNOWN  Status = "unknown"
)

// Service holds all information related to a service
type Service struct {
	name      string
	timeout   uint64
	time      chan uint64
	isHandled bool
}

var services = map[string]*Service{}

func main() {
	fmt.Println("Server listening on port 10000.")
	http.HandleFunc("/", handleRequests())
	log.Fatal(http.ListenAndServe(":10000", nil))
}

func handleRequests() func(w http.ResponseWriter, r *http.Request) {
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(fmt.Errorf("%+v", "Could not connect to docker API"))
	}
	return func(w http.ResponseWriter, r *http.Request) {
		serviceName, serviceTimeout, err := parseParams(r)
		if err != nil {
			fmt.Fprintf(w, "%+v", err)
		}
		service := GetOrCreateService(serviceName, serviceTimeout)
		status, err := service.HandleServiceState(cli)
		if err != nil {
			fmt.Printf("Error: %+v\n ", err)
			fmt.Fprintf(w, "%+v", err)
		}
		fmt.Fprintf(w, "%+s", status)
	}
}

func getParam(queryParams url.Values, paramName string) (string, error) {
	if queryParams[paramName] == nil {
		return "", fmt.Errorf("%s is required", paramName)
	}
	return queryParams[paramName][0], nil
}

func parseParams(r *http.Request) (string, uint64, error) {
	queryParams := r.URL.Query()

	serviceName, err := getParam(queryParams, "name")
	if err != nil {
		return "", 0, nil
	}

	timeoutString, err := getParam(queryParams, "timeout")
	if err != nil {
		return "", 0, nil
	}
	serviceTimeout, err := strconv.Atoi(timeoutString)
	if err != nil {
		return "", 0, fmt.Errorf("timeout should be an integer")
	}
	return serviceName, uint64(serviceTimeout), nil
}

// GetOrCreateService return an existing service or create one
func GetOrCreateService(name string, timeout uint64) *Service {
	if services[name] != nil {
		return services[name]
	}
	service := &Service{name, timeout, make(chan uint64), false}

	services[name] = service
	return service
}

// HandleServiceState up the service if down or set timeout for downing the service
func (service *Service) HandleServiceState(cli *client.Client) (string, error) {
	status, err := service.getStatus(cli)
	if err != nil {
		return "", err
	}
	if status == UP {
		fmt.Printf("- Service %v is up\n", service.name)
		if !service.isHandled {
			go service.stopAfterTimeout(cli)
		}
		select {
		case service.time <- service.timeout:
		default:
		}
		return "started", nil
	} else if status == STARTING {
		fmt.Printf("- Service %v is starting\n", service.name)
		if err != nil {
			return "", err
		}
		go service.stopAfterTimeout(cli)
		return "starting", nil
	} else if status == DOWN {
		fmt.Printf("- Service %v is down\n", service.name)
		service.start(cli)
		return "starting", nil
	} else {
		fmt.Printf("- Service %v status is unknown\n", service.name)
		if err != nil {
			return "", err
		}
		return service.HandleServiceState(cli)
	}
}

func (service *Service) getStatus(client *client.Client) (Status, error) {
	ctx := context.Background()
	dockerService, err := service.getDockerService(ctx, client)

	if err != nil {
		return "", err
	}

	if *dockerService.Spec.Mode.Replicated.Replicas == zeroReplica {
		return DOWN, nil
	}
	return UP, nil
}

func (service *Service) start(client *client.Client) {
	fmt.Printf("Starting service %s\n", service.name)
	service.isHandled = true
	service.setServiceReplicas(client, 1)
	go service.stopAfterTimeout(client)
	service.time <- service.timeout
}

func (service *Service) stopAfterTimeout(client *client.Client) {
	service.isHandled = true
	for {
		select {
		case timeout, ok := <-service.time:
			if ok {
				time.Sleep(time.Duration(timeout) * time.Second)
			} else {
				fmt.Println("That should not happen, but we never know ;)")
			}
		default:
			fmt.Printf("Stopping service %s\n", service.name)
			service.setServiceReplicas(client, 0)
			return
		}
	}
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
