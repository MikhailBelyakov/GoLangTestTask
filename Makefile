.PHONY: all

-include .env

SHELL=/bin/bash -e

.DEFAULT_GOAL := help

#####################
### Docker config ###
#####################

build-base:  ## Build docker image for nomad
	@docker-compose build


.PHONY: start
start: ## Run containers
	docker-compose up --build -d

.PHONY: stop
stop: ## Stop containers
	docker-compose down

.PHONY: ps
ps: ## List containers
	docker-compose ps

