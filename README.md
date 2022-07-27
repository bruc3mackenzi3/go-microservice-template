# microservice-demo
A microservice implementation for demonstration purposes

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

## ToDo
* Logging
* Context
* Integration tests
* Test line coverage report
* Interservice communication e.g. with gRPC
* (CICD)[https://docs.docker.com/language/golang/configure-ci-cd/]
* GitHub integrations
* Cloud deployment
