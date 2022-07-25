POSTGRES_CONTAINER = postgres-db
PGHOST = localhost
PGDB = microservice
PGUSER = microservice
PGPASSWORD = example123

build:
	docker-compose build

build-multistage:
	docker build --tag microservice-demo:multistage -f Dockerfile.multistage .

build-local:
	go build -v

start:
	docker-compose up --detach
	# docker cp ./sql/seed.sql $(POSTGRES_CONTAINER):/docker-entrypoint-initdb.d/seed.sql
	# docker exec -u postgres $(POSTGRES_CONTAINER) psql $(PGDB) $(PGUSER) -f docker-entrypoint-initdb.d/seed.sql

stop:
	docker-compose stop

run:
	PGHOST=$(PGHOST) PGDB=$(PGDB) PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) \
	./microservice-demo

clean:
	go clean ./...
	docker-compose rm --force

nuke:
	rm -rf vendor/
	docker rmi
	docker rmi microservice-demo postgres:14
