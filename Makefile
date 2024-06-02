DOCKER_COMPOSE_FILES ?= $(shell find docker -maxdepth 1 -type f -name "*.yaml" -exec printf -- '-f %s ' {} +; echo)


#######################################################################################################################

## ▸▸▸ Docker commands ◂◂◂

.PHONY: config
config:			## Show Docker config
	docker compose ${DOCKER_COMPOSE_FILES} config

.PHONY: up
up:			## Run Docker services
	docker compose ${DOCKER_COMPOSE_FILES} up --build --detach

.PHONY: down
down:			## Stop Docker services
	docker compose ${DOCKER_COMPOSE_FILES} down

.PHONY: ps
ps:			## Show Docker containers info
	docker ps --size --all --filter "name=gophermart-api"
