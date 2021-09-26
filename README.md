# Traefik on demand service

![Build](https://github.com/acouvreur/traefik-ondemand-service/workflows/Build/badge.svg)
![Latest version](https://img.shields.io/github/v/tag/acouvreur/traefik-ondemand-service?sort=semver)

## Description

This is a service that can scale up or down a docker swarm service on demand.
It basically starts a service when it's needed and then shut it down when it's no longer needed.

## Usage

### CLI

`./traefik-ondemand-service --swarmMode=true`

| Argument    | Value             | Description                                                             |
| ----------- | ----------------- | ----------------------------------------------------------------------- |
| `swarmMode` | true,false (default true) | Enable/Disable swarm mode. Used to determine the scaler implementation. |

### Docker

`docker run -v /var/run/docker.sock:/var/run/docker.sock -p 10000:10000 ghcr.io/acouvreur/traefik-ondemand-service:latest --swarmode=true`

### API

```
GET <service_url>:10000/?name=<service_name>&timeout=<timeout>
```

| Query param | Type            | Description                                                             |
| ----------- | --------------- | ----------------------------------------------------------------------- |
| `name`      | `string`        | The docker container name, or the swarm service name                    |
| `timeout`   | `time.Duration` | The duration after which the container/service will be scaled down to 0 |

| Body       | Status code  | Description                                                                    |
| ---------- | ------------ | ------------------------------------------------------------------------------ |
| `started`  | 202 Created  | The container/service is available                                             |
| `starting` | 201 Accepted | The container/service has been scheduled for starting but is not yet available |
