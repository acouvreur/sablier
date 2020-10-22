# Traefik Ondemand Plugin

Traefik plugin to scale down to zero containers on docker swarm.

## Description

A container may be a simple nginx server serving static pages, when no one is using it, it still consume CPU and memory.

With this plugin you can scale down to zero when there is no request for the service.
It will scale back to 1 when there is a user requesting the service.

## Demo

The service **whoami** is scaled to 0. We configured a **timeout of 10** seconds.

![Demo](./img/demo.gif)

## Run the demo

*use `watch -n 1 docker service ls` to see in real time the service getting downscaled*

1. `docker swarm init`
2. `export TRAEFIK_PILOT_TOKEN=your_traefik_pilot_token`
3. `docker stack deploy -c docker-compose.yml TRAEFIK_HACKATHON`
4. Go to `localhost:8000/whoami` --> service is starting
5. Refresh --> service is responding
6. wait 10 seconds
7. Refresh --> service is starting again because it was scaled down to 0

## Configuration

- `serviceUrl` the traefik-ondemand-service url (e.g. http://ondemand:1000)
- `name` the service to scale on demand name (docker service ls)
- *`timeout` (default: 60)* timeout in seconds for the service to be scaled down to zero after the last request

See `config.yml` and `docker-compose.yml` for full configuration.

## Limitations

### Cannot use service labels

You cannot set the labels for a service inside the service definition.

Otherwise when scaling to 0 the specification would not be found because there is no more task running. So you have to write it under the dynamic configuration file.

### The need of "traefik-ondemand-service"

We are running "traefik-ondemand-service" to interact freely with the docker deamon and manage an independant lifecycle from traefik.

*We may try to update this plugin to embed the scaling behavior in a future.*

-> The source is available at https://github.com/acouvreur/traefik-ondemand-service

## TODO

- [ ] Embed "traefik-ondemand-service" inside the plugin directly
- [ ] Scale **up** service (max replica, threshold)
- [ ] Scale down from N to 1 (threshold)
- [ ] Kubernetes integration
- [ ] Add configuration sample with plugin in non dev mode

## Authors

[Alexis Couvreur](https://www.linkedin.com/in/alexis-couvreur/) (left) and [Alexandre Hiltcher](https://www.linkedin.com/in/alexandre-hiltcher/) (right)

![Alexandre and Alexis](./img/gophers-traefik.png)
