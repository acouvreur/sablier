# Traefik Ondemand Plugin

Traefik plugin to scale down to zero containers on swarm.

## Description

A container may be a simple nginx server serving static pages, when no one is using it, it still consome CPU and memory even if it's close to nothing it still something.

With this plugin you can scale down to zero when there is no more request.
It will scale back to 1 when there is a user requesting the service.

## Authors

[Alexandre Hiltcher](https://www.linkedin.com/in/alexandre-hiltcher/)
[Alexis Couvreur](https://www.linkedin.com/in/alexis-couvreur/)

![Alexandre and Alexis](./img/gophers-traefik.png)

## Demo

![Demo](./img/demo.gif)

## Run the demo

- `docker swarm init`
- `export TRAEFIK_PILOT_TOKEN=your_traefik_pilot_token`
- `docker stack deploy -c docker-compose.yml TRAEFIK_HACKATHON`

## Limitations

You cannot set the labels for a service inside the service definition.

Otherwise when scaling to 0 the specification would not be found because there is no more task running. So you have to write it under the dynamic configuration file.