build:
	docker build --tag microservice-demo .

build-multistage:
	docker build --tag microservice-demo:multistage -f Dockerfile.multistage .

build-local:
	go build -v

start:
	docker run --publish 8080:8080 --detatch microservice-demo --name microserver

stop:
	docker stop microserver

run:
	./microservice-demo

clean:
	go clean ./...
	docker rm microserver
	docker rmi microservice-demo
