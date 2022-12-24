POSTGRES_CONTAINER = postgres-db
PGHOST = localhost
PGDB = microservice
PGUSER = microservice
PGPASSWORD = example123

# help extracts the help texts for the comments following ': ##'
.PHONY: help
help: ## Print this help message
	@awk -F':.*## ' ' \
		/^[[:alpha:]_-]+:.*## / { \
			printf "\033[36m%s\033[0m\t%s\n", $$1, $$2 \
		} \
	' $(MAKEFILE_LIST) | column -s$$'\t' -t

.PHONY: init
init:
	go mod tidy
	go mod vendor
	go install github.com/vektra/mockery/v2@latest

.PHONY: build
build:  ## Compile Go program
	go build -v

.PHONY: build-multistage
build-multistage:  ## Minimal Docker build resulting in a much smaller container
	docker build --tag microservice-demo:multistage -f Dockerfile.multistage .

.PHONY: start-platform
start-platform:  ## Build & start platform container(s)
	docker-compose build
	docker-compose --profile platform up --detach
	# These commands seed the database; currently disabled in favour of GORM AutoMigrate
	# docker cp ./sql/seed.sql $(POSTGRES_CONTAINER):/docker-entrypoint-initdb.d/seed.sql
	# docker exec -u postgres $(POSTGRES_CONTAINER) psql $(PGDB) $(PGUSER) -f docker-entrypoint-initdb.d/seed.sql

.PHONY: run
run:  ## Run application natively
	PGHOST=$(PGHOST) PGDB=$(PGDB) PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) \
	go run main.go

.PHONY: start
start:  ## Start application container(s)
	# WARNING: make build does not correctly rebuild users container when the image already exists
	docker-compose stop debug-app
	docker-compose --profile application up --detach

.PHONY: stop
stop:  ## Stop all containers
	docker-compose --profile application --profile platform stop

.PHONY: start-debug
start-debug:  ## Restart app in debug mode with Delve debugger
	docker-compose stop app  # Note: must specify service directly since profiles doesn't work with stop
	docker-compose --profile debug build
	docker-compose --profile debug up --detach

.PHONY: test
test:  ## Run test suite
	go test -v ./...

.PHONY: clean
clean:  ## Clean project
	go clean ./...
	docker-compose rm --force

.PHONY: nuke
nuke:  ## Remove all data generated with the project
	rm -rf vendor/
	docker rmi microservice-demo postgres:14
	docker network rm microservice_network
