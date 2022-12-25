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
* Docker Compose for container orchestration
* VSCode Integration
  * launch.json configurations for:
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

## Developing
### Mocks
Update a mock with the following command, e.g. for the `repository.Repository` interface:
```bash
cd repository/
mockery --name Repository --inpackage --outpkg=mock_Repository.go
```

## ToDo
* BUG: Docker app cannot connect to DB - likely needs a Docker network configured
* Graceful shutdown (should entail catching SIGTERM, disconnecting from DB & exiting)
* Logging
* Context
* Integration tests
* Test line coverage report
* Interservice communication e.g. with gRPC
* (CICD)[https://docs.docker.com/language/golang/configure-ci-cd/]
* GitHub integrations
* Cloud deployment
