# go-microservice-template

A full-featured, production-ready microservice template for building applications in Go.  Includes a Users microservice exposing a REST API providing functionality to perform CRUD operations on the resource.

[Fork this project](https://github.com/bruc3mackenzi3/go-microservice-template/fork) to start building your own REST-based microservice application.

![example workflow](https://github.com/bruc3mackenzi3/microservice-demo/actions/workflows/go.yml/badge.svg)

## Tech Stack
* Makefile for ease of operation
  * Encapsulates commands to build, run, test, clean, etc, the project
* Go
  * [Echo web framework](https://github.com/labstack/echo)
  * [GORM ORM Library](https://github.com/go-gorm/gorm)
  * [Testify](https://github.com/stretchr/testify) Assert and Mock packages
  * golangci-lint with `.golangci.yml` config and VSCode `settings.json` integration
* PostgreSQL 14 database
* Docker for containerizing the application
* Kubernetes for deploying the containerized application
* VSCode Integration
  * `launch.json` configurations for:
    * Debugging native application
    * Debugging in Docker container using [Delve debugger](https://github.com/go-delve/delve)
* GitHub Actions to run build & tests
* Docs
  * README
  * OpenAPI definition documenting REST API

## Getting Started
Fetch Go dependencies
```
make init
```

Compile the Go application
```
make build
```

Start the database
```
make start-db
```

Run the microservice app directly
```
make run
```

Execute unit tests
```
make test
```

Stop the database
```
make stop-db
```

Get help with all make commands
```
make help
```

Query the application using Curl:
```bash
curl -X POST http://localhost/users/Bruce
curl http://localhost/users/1
```

### Running with Docker
Run the microservice app in Docker
```sh
make d-build
make d-start
```

stop the app
```sh
make stop
```

Rebuild
```bash
docker rm users
make d-build
make d-start
```

### Running in Kubernetes
Precondition: kubectl is installed in your local environment, and connected to a Kubernetes cluster.

First deploy the database:
```sh
kubectl apply -f service.postgres.yaml
```

Next deploy the Users service:
```sh
kubectl apply -f service.users.yaml
```

Tail application logs:
```
kubectl logs -f users
```

The service is exposed externally with a Kubernetes object called a NodePort Service.  This NodePort is exposed externally on port `30080`.  If you're using a local K8s cluster like MiniKube or Docker Desktop, it will be accessible via `localhost`, or `127.0.0.1`:
```
curl http://localhost:30080/v1/users/1
```

The Postgres database is exposed internally only, because in a real world scenario it should not be accessible from outside the Kubernetes cluster.  Assuming your cluster is running locally, you can connect to it by telling Kubernetes to expose the connection:
```
kubectl port-forward postgres 5432
```
Then connect to it as you normally would, e.g.:
```
psql -U microservice
```

Delete application resources created on the cluster:
```sh
kubectl delete service postgres users
kubectl delete pod postgres users
```

Here is a diagram of the setup in Kuberenetes:
![alt text](./k8s.png)

## Developing
### Kubernetes
To run your changes in Kubernetes follow these steps:

1. Log in to Docker
```
docker login -U username
```

2. Set `DOCKER_USERNAME` in `Makefile` to your Docker username

3. Build and tag the Docker image
```
make d-build
```

4. Push the Docker image
```
make d-push
```

5. Deploy to Kubernetes with the instructions under "Running in Kubernetes" above

### Mocks
Update a mock with the following command, e.g. for the `repository.Repository` interface:
```bash
cd repository/
mockery --name Repository --inpackage --outpkg=mock_Repository.go
```

## ToDo
* Host OpenAPI @ github.io
* BUG: Docker app cannot connect to DB - likely needs a Docker network configured
* Logging
* Context
* Integration tests
* Test line coverage report
* Inter-service communication e.g. with gRPC
* (CICD)[https://docs.docker.com/language/golang/configure-ci-cd/]
* GitHub integrations
* Cloud deployment
