# Traefik Ondemand Plugin

Traefik plugin to scale containers ondemand

## How to use it

### Required

- Swarm mode
- [Traefik ondemand service](https://github.com/acouvreur/traefik-ondemand-service) up and running

## Develop the plugin

`export TRAEFIK_PILOT_TOKEN=traefik_pilot_token`
`docker swarm init`
`docker stack deploy -c docker-compose.yml TRAEFIK_HACKATHON`