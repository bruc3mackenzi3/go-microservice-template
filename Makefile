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

start-platform:
	docker-compose --profile platform up --detach

start:
	docker-compose --profile application up --detach
	# These commands seed the database; currently disabled in favour of GORM AutoMigrate
	# docker cp ./sql/seed.sql $(POSTGRES_CONTAINER):/docker-entrypoint-initdb.d/seed.sql
	# docker exec -u postgres $(POSTGRES_CONTAINER) psql $(PGDB) $(PGUSER) -f docker-entrypoint-initdb.d/seed.sql

start-debug:
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

stop:
	docker-compose --profile application stop

stop-platform:
	docker-compose --profile platform stop

run:
	PGHOST=$(PGHOST) PGDB=$(PGDB) PGUSER=$(PGUSER) PGPASSWORD=$(PGPASSWORD) \
	./microservice-demo

clean:
	go clean ./...
	docker-compose rm --force
	docker rm --force microserver-debug

nuke:
	rm -rf vendor/
	docker rmi microservice-demo postgres:14
	docker network rm microservice_network
