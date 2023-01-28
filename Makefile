# Set go executable
# override this variable with an environment variable to run against another version of Go, e.g.:
# make GO=go1.19.4 build
GO = go

DOCKER_USERNAME := bruc3mackenzi3

PG_IMAGE = postgres:14
PG_CONTAINER = postgres-db
USERS_IMAGE = users
USERS_CONTAINER = users
USERS_PORT = 8080
USERS_VERSION := 1.0

PG_HOST = localhost
PG_DB = microservice
PG_USER = microservice
PG_PASSWORD = example123

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
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) install github.com/vektra/mockery/v2@latest
	docker pull $(PG_IMAGE)
	docker create --name $(PG_CONTAINER) --publish 5432:5432 -e POSTGRES_DB=$(PG_DB) -e POSTGRES_USER=$(PG_USER) -e POSTGRES_PASSWORD=$(PG_PASSWORD) $(PG_IMAGE)

.PHONY: build
build:  ## Compile Go program
	$(GO) build -v -o users

.PHONY: start-db
start-db:  ## Start the Postgres Database
	docker start $(PG_CONTAINER)

.PHONY: run
run:  ## Run the Users application (exit with CTRL+C)
	POSTGRES_DB=$(PG_DB) POSTGRES_USER=$(PG_USER) POSTGRES_PASSWORD=$(PG_PASSWORD) POSTGRES_HOST=$(PG_HOST) \
	$(GO) run main.go

.PHONY: test
test:  ## Run the unit test suite
	$(GO) test -v ./...

.PHONY: integration-test
integration-test:  ## Run the integration test suite
	cd integration_tests/; \
	curl --silent http://localhost:8080/healthz > /dev/null && \
	POSTGRES_DB=$(PG_DB) POSTGRES_USER=$(PG_USER) POSTGRES_PASSWORD=$(PG_PASSWORD) POSTGRES_HOST=$(PG_HOST) \
	INTEGRATION_TESTS=1 go test . || echo "Failed to reach service; cannot run test"

.PHONY: stop-db
stop-db:
	docker stop $(PG_CONTAINER)

.PHONY: d-build
d-build:  ## Build app in Docker image
	env GOOS=linux GOARCH=amd64 $(GO) build -o users-docker
	docker build -t $(USERS_IMAGE):latest .
	docker tag $(USERS_IMAGE) $(DOCKER_USERNAME)/$(USERS_IMAGE):$(USERS_VERSION)

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

.PHONY: d-push
d-push:  ## Push built Users image to public Docker repository.  Must be logged in to Docker.
	docker push $(DOCKER_USERNAME)/$(USERS_IMAGE):$(USERS_VERSION)

.PHONY: clean
clean:  ## Clean project
	$(GO) clean ./...
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
