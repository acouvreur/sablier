# Traefik Ondemand Plugin

![Build](https://github.com/acouvreur/traefik-ondemand-plugin/workflows/Build/badge.svg)

Start your containers/services on the first request they recieve, and shut them down after a specified duration after the last request they received. 

Docker classic and docker swarm compatible.

![Demo](./img/ondemand.gif)
## Usage

```yml
testData:
  serviceUrl: http://ondemand:10000
  name: TRAEFIK_HACKATHON_whoami
  timeout: 1m
```

| Parameter    | Type            | Example                    | Description                                                             |
| ------------ | --------------- | -------------------------- | ----------------------------------------------------------------------- |
| `serviceUrl` | `string`        | `http://ondemand:10000`    | The docker container name, or the swarm service name                    |
| `name`       | `string`        | `TRAEFIK_HACKATHON_whoami` | The container/service to be stopped (docker ps                          | docker service ls) |
| `timeout`    | `time.Duration` | `1m30s`                    | The duration after which the container/service will be scaled down to 0 |


## Examples

- [Docker Classic](./examples/docker_classic/)
- [Docker Swarm](./examples/docker_swarm/)

## Authors

[Alexis Couvreur](https://www.linkedin.com/in/alexis-couvreur/) (left) and [Alexandre Hiltcher](https://www.linkedin.com/in/alexandre-hiltcher/) (right)

![Alexandre and Alexis](./img/gophers-traefik.png)
