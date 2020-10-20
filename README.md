# Traefik Ondemand Plugin

Traefik plugin to scale containers ondemand

## Required

- Service to scale is a swarm service

## Develop the plugin

`export TRAEFIK_PILOT_TOKEN=traefik_pilot_token`
`docker swarm init`
`docker stack deploy -c docker-compose.yml TRAEFIK_HACKATHON`