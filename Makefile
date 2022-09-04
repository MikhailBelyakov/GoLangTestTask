.PHONY: all

-include .env

SHELL=/bin/bash -e

.DEFAULT_GOAL := help

#####################
### Docker config ###
#####################

.PHONY: help
help: ## Show this message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build:  ## Build docker image
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

