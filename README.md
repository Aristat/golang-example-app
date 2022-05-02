![golang logo](golang_logo.png)

# Golang Example Application

# Table of Contents

- [Overview](#overview)
- [Package list](#package-list)
- [Installing](#installing)
  * [Local environment](#local-environment)
  * [Docker environment](#docker-enviroment)
- [Run services](#run-services)
  * [Start in local](#start-in-local-machine)
  * [Start in docker](#start-in-docker)
- [Getting started](#getting-started)  
  * [Jaeger](#jaeger)
  * [Http with gRPC](#http-example-with-grpc)
  * [Graphql with gRPC](#graphql-example-with-grpc)
- [Deprecate version](#deprecated-version)  
- [Testing](#testing)

# Overview

This is an example golang application.
Commands list:
1. Daemon - main service
2. Product service - service that returns Product, an example of gRPC client/server interaction
3. Health check service - this service is needed to show how convenient to understand on which of the services an error occurred in jaeger
4. Migrate - commands for migration
5. JWT - commands for generate JWT token

# Package list

Packages which use in this example project

1. [sql-migrate](https://github.com/rubenv/sql-migrate) - SQL migrations
2. [wire](https://github.com/google/wire) - dependency Injection
3. [viper](https://github.com/spf13/viper) - environment configuration
4. [cobra](https://github.com/spf13/cobra) - create commands
5. [cast](https://github.com/spf13/cast) - easy casting from one type to another
6. [gorm](https://github.com/jinzhu/gorm) - database ORM
7. [zap](https://github.com/uber-go/zap) - logger
8. [mux](https://github.com/gorilla/mux) - http router
9. [nats-streaming](https://github.com/nats-io/stan.go) - NATS Streaming System
10. [gqlgen](https://github.com/99designs/gqlgen) - graphql server library
11. [protobuf](https://pkg.go.dev/google.golang.org/protobuf) - Google's data interchange format
12. [grpc](google.golang.org/grpc) - RPC framework
13. [opentelemetry](https://github.com/open-telemetry/opentelemetry-go) - OpenTelemetry
14. [jaeger](https://github.com/uber/jaeger-client-go) - Jaeger Bindings for Go OpenTelemetry API
15. [casbin](https://github.com/casbin/casbin) - Supports access control
16. [dataloaden](https://github.com/vektah/dataloaden) - DataLoader for graphql
17. [nats](https://github.com/nats-io/nats.go) - Golang client for NATS, the cloud native messaging system

# Installing

Install the Golang and GO environment

```$xslt
https://golang.org/doc/install
```

Install [Postgresql](https://www.postgresql.org/download) (if you want to run locally)

Clone repository

```$xslt
git clone git@github.com:Aristat/golang-example-app.git (go get)
```

## Local environment

Install Golang packages without modules

```$xslt
make install
```

Install database

```$xslt
make createdb
```

Sql migrations

```$xslt
sql-migrate up
```

Install Golang dependencies

```$xslt
make dependencies
```

Generate artifacts(binary files and configs)

```$xslt
make build
```

Packages for proto generator

```$xslt
https://grpc.io/docs/languages/go/quickstart/#prerequisites
```

Set APP_WD if you start to use html templates or path to ssh keys or run `make test`

```$xslt
export APP_WD=go_path to project_path/resources or project_path/artifacts
```

## Docker environment

Generate docker image

```$xslt
DOCKER_IMAGE=golang-example-app TAG=development make docker-image
```

# Run services

## Start in local machine

Up jaeger in docker-compose or disable Jaeger(and rebuild binary file) in `resources/configs/*.yaml`

```$xslt
docker-compose up jaeger
```

Start daemon (main service)

```$xslt
make start
```

or

```$xslt
./artifacts/bin daemon -c ./artifacts/configs/development.yaml -d
```

Start product service 

```$xslt
./artifacts/bin product-service -c ./artifacts/configs/development.yaml -d
```

Start health-check service

```$xslt
./artifacts/bin health-check -c ./artifacts/configs/development.yaml -d
```

## Start in docker

#### Run this commands
```$xslt
docker-compose rm # Remove previous containers
REMOVE_CONTAINERS=on DOCKER_IMAGE=golang-example-app TAG=development make docker-image # Generate new docker image
docker-compose up
```

#### or run script

```$xslt
./scripts/docker-compose-start.sh
```

# Getting Started

## Jaeger

```$xslt
http://localhost:16686
```

## Http example with gRPC

```$xslt
http://localhost:9096/products_grpc
```

## Http example with Nats Streaming

```$xslt
http://localhost:9096/products_nats
```

## Http example with graceful shutdown(long request, you can check the server shutdown)

```$xslt
http://localhost:9096/products_slowly
```

## Graphql example with gRPC

Graphql [client](https://github.com/prisma-labs/graphql-playground) for testing. End-point `http://localhost:9096/query`.

Generate JWT for Graphql authorization

```$xslt
./artifacts/bin jwt token --key='./artifacts/keys/local/private_key.pem' --fields='{"sub": "owner", "iss": "test-service"}'
```

Set JWT token in headers

```$xslt
{
  "Authorization": "bearer token"
}
```

Example query

```
query oneUser {
  users {
    one(email: "test@gmail.com") {
      email
      id
    }
  }
}
```

Example query with data loader

```
query allProducts {
  products {
    list {
      products {
        id
        productItems {
          id
          name
        }
      }
    }
  }
}
```

Example mutation

```
mutation createUser {
  users {
    createUser(email: "test1@gmail.com", password: "123456789") {
      id
      email
    }
  }
}
```

# Testing
```
➜  make test
```
