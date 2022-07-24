build:
	docker build --tag microservice-demo .

build-multistage:
	docker build --tag microservice-demo:multistage -f Dockerfile.multistage .

build-local:
	go build -v

start:
	docker run --publish 8080:8080 -d microservice-demo

run:
	./microservice-demo

clean:
	go clean ./...
