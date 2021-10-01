# Docker swarm

1. `docker swarm init`
2. `docker stack deploy -c docker-stack.yml DOCKER_SWARM`
3. Load `http://localhost/nginx`
4. Wait 1 minute
5. Service is scaled to 0/0