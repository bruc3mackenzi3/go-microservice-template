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

build:  ## Compile Go program
	go build -v

build-multistage:  ## Minimal Docker build resulting in a much smaller container
	docker build --tag microservice-demo:multistage -f Dockerfile.multistage .

start-platform:  ## Build & start platform container(s)
	docker-compose build
	docker-compose --profile platform up --detach

run:  ## Run application natively
	PGHOST=$(PGHOST) PGDB=$(PGDB) PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) \
	go run main.go

start:  ## Start application container(s)
	docker-compose --profile application up --detach
	# These commands seed the database; currently disabled in favour of GORM AutoMigrate
	# docker cp ./sql/seed.sql $(POSTGRES_CONTAINER):/docker-entrypoint-initdb.d/seed.sql
	# docker exec -u postgres $(POSTGRES_CONTAINER) psql $(PGDB) $(PGUSER) -f docker-entrypoint-initdb.d/seed.sql

start-debug:  ## Restart app in debug mode with Delve debugger
	docker-compose --profile application stop
	docker rm --force microserver-debug
	docker build -f Dockerfile.debug --tag debug .
	docker run --name microserver-debug \
		--network microservice_network \
		--env PGHOST=postgres-db \
		--env PGDB=microservice \
		--env PGUSER=microservice \
		--env PGPASSWORD=example123 \
		--publish 8080:80 \
		--publish 4000:4000 \
		debug

stop:  ## Stop application container(s)
	docker-compose --profile application stop

stop-platform:  ## Stop platform container(s)
	docker-compose --profile platform stop

test:  ## Run test suite
	go test -v ./...

clean:  ## Clean project
	go clean ./...
	docker-compose rm --force
	docker rm --force microserver-debug

nuke:  ## Remove all data generated with the project
	rm -rf vendor/
	docker rmi microservice-demo postgres:14
	docker network rm microservice_network
