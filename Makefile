DOCKER_COMPOSE_FILENAME=deployments/docker-compose.yml

start:
	docker-compose --file ${DOCKER_COMPOSE_FILENAME} up -d
