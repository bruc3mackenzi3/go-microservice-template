SHELL := /bin/bash
PG_HOST = localhost
PG_DB = microservice
PG_USER = microservice
PG_PASSWORD = example123

PG_IMAGE = postgres:14
USERS_IMAGE = users
PG_CONTAINER = postgres-db
USERS_CONTAINER = users
USERS_PORT = 8080

# help extracts the help texts for the comments following ': ##'
.PHONY: help
help: ## Print this help message
	@awk -F':.*## ' ' \
		/^[[:alpha:]_-]+:.*## / { \
			printf "\033[36m%s\033[0m\t%s\n", $$1, $$2 \
		} \
	' $(MAKEFILE_LIST) | column -s$$'\t' -t

.PHONY: init
init:  ## Fetch Go dependencies; Download Postgres Docker Image
	go mod tidy
	go mod vendor
	go install github.com/vektra/mockery/v2@latest
	docker pull $(PG_IMAGE)
	docker create --name $(PG_CONTAINER) --publish 5432:5432 -e POSTGRES_DB=$(PG_DB) -e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASSWORD) $(PG_IMAGE)

.PHONY: build
build:  ## Compile Go program
	go build -v

.PHONY: start-db
start-db:  ## Start the Postgres Database
	docker start $(PG_CONTAINER)

.PHONY: run
run:  ## Run the Users application (exit with CTRL+C)
	POSTGRES_DB=$(PG_DB) POSTGRES_USER=$(PG_USER) POSTGRES_PASSWORD=$(PG_PASSWORD) POSTGRES_HOST=$(PG_HOST) \
	go run main.go

.PHONY: stop-db
stop-db:
	docker stop $(PG_CONTAINER)

.PHONY: d-build
d-build:  ## Build app in Docker image
	env GOOS=linux GOARCH=amd64 go build -o users-docker
	docker build -t $(USERS_IMAGE):latest .

.PHONY: d-run
d-run:  ## Run the Users application in a Docker Container; requires previous containers are stopped
	docker create --name $(USERS_CONTAINER) --publish $(USERS_PORT):$(USERS_PORT) -e POSTGRES_DB=$(PG_DB) -e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASSWORD) -e POSTGRES_HOST=$(PG_HOST) $(USERS_IMAGE)
	docker start $(USERS_CONTAINER)

.PHONY: d-stop
d-stop:  ## Stop Users container
	docker stop $(USERS_CONTAINER) users-debug
	docker rm $(USERS_CONTAINER)

.PHONY: d-start-debug
start-debug:  ## Restart app in debug mode with Delve debugger
	docker build -t users-debug -f Dockerfile.debug
	docker start users-debug

.PHONY: test
test:  ## Run test suite
	go test -v ./...

.PHONY: clean
clean:  ## Clean project
	go clean ./...
	docker rm --force postgres-db users

.PHONY: nuke
nuke:  ## Remove all data generated with the project
	rm -rf vendor/
	docker rm --force $(PG_CONTAINER) $(USERS_CONTAINER)
	docker rmi $(USERS_IMAGE) $(PG_IMAGE)

.PHONY: d-build-multistage
d-build-multistage:  ## Minimal Docker build resulting in a much smaller container
	docker build --tag users:multistage -f Dockerfile.multistage .
	docker create --name users-multistage --publish $(USERS_PORT):$(USERS_PORT) -e POSTGRES_DB=$(PG_DB) -e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASSWORD) -e POSTGRES_HOST=$(PG_HOST) users:multistage
