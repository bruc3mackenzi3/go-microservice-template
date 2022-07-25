build:
	docker-compose build

build-multistage:
	docker build --tag microservice-demo:multistage -f Dockerfile.multistage .

build-local:
	go build -v

start:
	docker-compose up --detatch

stop:
	docker-compose down

run:
	./microservice-demo

clean:
	go clean ./...
	# docker-compose rm --force

nuke:
	rm -rf vendor/
	docker rmi
	docker rmi microservice-demo postgres:14
